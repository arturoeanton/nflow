package playbook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/arturoeanton/gocommons/utils"
	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/dop251/goja_nodejs/util"

	"github.com/labstack/echo/v4"
)

var (
	registry *require.Registry
	jsVars   map[string]string = make(map[string]string)
)

func (cc *Controller) Run(c echo.Context, vars Vars, next string, uuid1 string, payload goja.Value) error {
	return cc.run(c, vars, next, uuid1, payload, false)
}

func (cc *Controller) RunWithCallback(c echo.Context, vars Vars, next string, uuid1 string, payload goja.Value) error {
	return cc.run(c, vars, next, uuid1, payload, true)
}

func (cc *Controller) run(c echo.Context, vars Vars, next string, uuid1 string, payload goja.Value, fork bool) error {
	pathBase := GetPathBase(c)

	var p *process.Process

	if fork {
		p = process.CreateProcessWithCallback(uuid1)
		go func(uuid2 string, currentProcess *process.Process) {
			data := <-currentProcess.Callback
			var p map[string]interface{}
			json.Unmarshal([]byte(data), &p)
			if _, ok := p["error_exit"]; ok {
				currentProcess.FlagExit = 1
			}

		}(uuid1, p)
	} else {
		p = process.CreateProcess(uuid1)
	}

	defer func() {
		p.SendCallback(`{"error_exit":"exit"}`)
		p.Close()
	}()

	if c.Response().Header().Get("Nflow-Wid-1") == "" {
		c.Response().Header().Add("Nflow-Wid-1", uuid1)
	}

	vm := goja.New()

	if registry == nil {
		registry = new(require.Registry) // this can be shared by multiple runtimes
		registry.RegisterNativeModule("console", console.Require)
		registry.RegisterNativeModule("util", util.Require)
	}

	registry.Enable(vm)
	console.Enable(vm)

	addFeatureSession(vm, c)
	addFeatureWsConsole(vm, c)

	addGlobals(vm, c)

	for _, p := range Plugins {
		for key, fx := range p.AddFeatureJS() {
			vm.Set(key, fx)
		}
	}

	vm.Set("c", c)
	vm.Set("echo_context", c)

	postData := make(map[string]interface{})
	func() {
		c.Bind(&postData)
		vm.Set("post_data", postData)
	}()

	vm.Set("vars", vars)
	vm.Set("path_vars", vars)

	vm.Set("wid", uuid1)

	vm.Set("wkill", func(wid string) {
		process.WKill(wid)
	})

	pb := *cc.Playbook
	node_auth := pb[next]
	if next == "" {
		next = cc.Start.Outputs["output_1"].Connections[0].Node
		node_auth = cc.Start
	}

	// Exceute AUTH.js?
	if flag, ok := node_auth.Data["nflow_auth"]; ok {
		flagString, ok := flag.(string)
		if !ok {
			flagBool := flag.(bool)
			flagString = fmt.Sprint(flagBool)
		}
		if flagString != "false" {
			//execute auth.js
			//nsess_auth, _ := app.Store.Get(c.Request(), "auth-session")
			//sess_auth.Values["redirect_url"] = c.Request().URL.Path
			///sess_auth.Save(c.Request(), c.Response())
			//profile := sess_auth.Values["profile"]
			//vm.Set("profile", profile)
			vm.Set("next", next)
			vm.Set("auth_flag", flagString)
			vm.Set("url_access", c.Request().URL.Path)

			pathAuth := pathBase + "auth.js"
			fmt.Println(pathAuth)
			code, _ := utils.FileToString(pathAuth)
			code += "\nmain()"
			_, err := vm.RunString(code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
			}

			next = vm.Get("next").String()
			fmt.Println(next)
			if next == "login" {
				return c.Redirect(http.StatusTemporaryRedirect, "/nflow_login")
			}
			if next == "break" {
				return nil
			}
		}
	}

	cc.Execute(c, vm, next, vars, p, payload)

	return nil
}

func (cc *Controller) step(c echo.Context, vm *goja.Runtime, next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error) {
	t1 := time.Now()
	sbLog := strings.Builder{}
	connection_next := "output_1"
	SendConsoleWS("__node_id:run:" + next)
	defer func() {
		diff := time.Now().Sub(t1)
		SendConsoleWS("__node_id:stop:" + next + ":" + fmt.Sprint(diff))
		log.Println(sbLog.String() + " - time:" + fmt.Sprint(diff))
	}()
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
			log.Println("step_00010 ****", err)
		}
	}()

	if currentProcess.FlagExit == 1 {
		currentProcess.Close()
		panic("FlagExit")
	}
	if false {
		log.Println(next)
	}
	pb := *cc.Playbook
	actor := pb[next]
	sbLog.WriteString("- IDBox:" + next)
	currentProcess.UUIDBoxCurrent = next

	if nameBox, ok := actor.Data["name_box"]; ok {
		sbLog.WriteString("- NameBox:" + nameBox.(string))
	}

	currentProcess.Type = ""
	if pType, ok := actor.Data["type"]; ok {
		currentProcess.Type = pType.(string)
	}

	sbLog.WriteString(" - Type:" + currentProcess.Type)
	if s, ok := Steps[currentProcess.Type]; ok {
		var err error
		connection_next, payload, err = s.Run(cc, actor, c, vm, connection_next, vars, currentProcess, payload)
		if err != nil {
			sbLog.WriteString(" - Error: " + err.Error())
			return "", nil, nil
		}
	} else {

		if currentProcess.Type == "starter" {
			c.JSON(http.StatusInternalServerError, echo.Map{"error": "Starter can not run with play button"})
			sbLog.WriteString(" - Error: Starter can not run with play button")
			return "", nil, nil
		}

		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Type node not found", "type": currentProcess.Type})
		sbLog.WriteString(" - Error: Not Found type")
		return "", nil, nil
	}

	sbLog.WriteString(" - Next:" + connection_next)
	return connection_next, payload, nil
}

func (cc *Controller) Execute(c echo.Context, vm *goja.Runtime, next string, vars Vars, currentProcess *process.Process, payload goja.Value) {
	var err error
	prev_box := ""
	for next != "" {

		vm.Set("current_box", next)
		vm.Set("prev_box", prev_box)
		prev_box = next
		next, payload, err = cc.step(c, vm, next, vars, currentProcess, payload)

		if err != nil {
			break
		}

		// cut
		if payload != nil {
			if rawPayload, ok := payload.Export().(map[string]interface{}); ok {
				if raw, ok := rawPayload["break"]; ok {
					if flag, ok := raw.(bool); ok {
						if flag {
							break
						}
					}
					if flag, ok := raw.(string); ok {
						if flag == "true" {
							break
						}
					}
				}
			}
		}

	}

}
