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

	"github.com/labstack/echo-contrib/session"
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

	// Exceute auth of default.js?
	if flag, ok := node_auth.Data["nflow_auth"]; ok {
		flagString, ok := flag.(string)
		if !ok {
			flagBool := flag.(bool)
			flagString = fmt.Sprint(flagBool)
		}
		if flagString != "false" {
			//execute auth of default.js
			auth_session, err := session.Get("auth-session", c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
				return nil
			}

			auth_session.Values["redirect_url"] = c.Request().URL.Path
			auth_session.Save(c.Request(), c.Response())

			profile := auth_session.Values["profile"]
			vm.Set("profile", profile)
			vm.Set("next", next)
			vm.Set("auth_flag", flagString)
			vm.Set("url_access", c.Request().URL.Path)

			defaultjs := pathBase + "default.js"
			code, err := utils.FileToString(defaultjs)
			if err != nil {
				c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
				return nil
			}

			code += "\nauth()"
			_, err = vm.RunString(code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
				return nil
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

	var actor *Node
	var box_id string
	var box_name string
	var box_type string
	defer func() {
		now := time.Now()
		diff := now.Sub(t1)

		go func(c echo.Context, actor *Node, box_id string, box_name string, box_type string, connection_next string, diff time.Duration) {
			//log.Println(sbLog.String() + " - time:" + fmt.Sprint(diff))
			pathBase := GetPathBase(c)
			defaultjs := pathBase + "default.js"
			code, err := utils.FileToString(defaultjs)
			if err != nil {
				log.Println(err)
				return
			}

			vm.Set("box_id", box_id)
			vm.Set("box_name", box_name)
			vm.Set("box_type", box_type)
			vm.Set("connection_next", connection_next)

			vm.Set("duration_mc", diff.Microseconds())
			vm.Set("duration_ms", diff.Milliseconds())
			vm.Set("duration_s", diff.Seconds())

			code += "\nlog()"
			_, err = vm.RunString(code)
			if err != nil {
				log.Println(err)
				return
			}

		}(c, actor, box_id, box_name, box_type, connection_next, diff)

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
	actor = pb[next]
	sbLog.WriteString("- IDBox:" + next)
	currentProcess.UUIDBoxCurrent = next
	box_id = next

	if nameBox, ok := actor.Data["name_box"]; ok {
		box_name = nameBox.(string)
		sbLog.WriteString("- NameBox:" + box_name)
	}

	currentProcess.Type = ""
	if pType, ok := actor.Data["type"]; ok {
		currentProcess.Type = pType.(string)
	}
	box_type = currentProcess.Type

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

		vm.Set("payload", payload)
		payload, _ := vm.RunString(`
		data = open_session("nflow_form")
		payload = {
			...payload,
			...data
		};
		payload;
		`)
		next, payload, err = cc.step(c, vm, next, vars, currentProcess, payload)
		if err != nil {
			break
		}

		// cut
		if payload != nil {
			if rawPayload, ok := payload.Export().(map[string]interface{}); ok {
				s, _ := session.Get("nflow_form", c)
				for k, v := range rawPayload {
					s.Values[k] = v
				}
				s.Save(c.Request(), c.Response())

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

	if next == "" {
		s, _ := session.Get("nflow_form", c)
		s.Values = make(map[interface{}]interface{})
		s.Save(c.Request(), c.Response())

		currentProcess.State = "end"
		currentProcess.Killeable = false
		currentProcess.Close()

	}

}
