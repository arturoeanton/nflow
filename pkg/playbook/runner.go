package playbook

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/dop251/goja_nodejs/util"
	"github.com/google/uuid"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	registry *require.Registry
	jsVars   map[string]string = make(map[string]string)
	wg       sync.WaitGroup    = sync.WaitGroup{}
)

func (cc *Controller) Run(c echo.Context, vars Vars, next string, uuid1 string, payload goja.Value) error {
	return cc.run(c, vars, next, uuid1, payload, false)
}

func (cc *Controller) RunWithCallback(c echo.Context, vars Vars, next string, uuid1 string, payload goja.Value) error {
	return cc.run(c, vars, next, uuid1, payload, true)
}

func (cc *Controller) run(c echo.Context, vars Vars, next string, uuid1 string, payload goja.Value, fork bool) error {

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
	addFeatureUsers(vm, c)
	addFeatureToken(vm, c)
	addFeatureTemplte(vm, c)

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

			ctx := c.Request().Context()
			db, err := GetDB()
			if err != nil {
				log.Println(err)
				return nil
			}
			conn, err := db.Conn(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
				return nil
			}
			defer conn.Close()
			row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetApp, "app")
			var code string
			var json_code string
			err = row.Scan(&json_code, &code)
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

	cc.Execute(c, vm, next, vars, p, payload, fork)

	return nil
}

func (cc *Controller) step(c echo.Context, vm *goja.Runtime, next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error) {
	t1 := time.Now()
	sbLog := strings.Builder{}
	connection_next := "output_1"

	log_session, err := session.Get("log-session", c)
	if err != nil {
		log.Println(err)
	}
	if log_session.Values["log_id"] == nil {
		log_session.Values["log_id"] = uuid.New().String()
		log_session.Values["order_box"] = 0
	}

	var log_id string
	var order_box int

	func() {

		log_id = log_session.Values["log_id"].(string)
		order_box = log_session.Values["order_box"].(int) + 1
		log_session.Values["order_box"] = order_box
		log_session.Save(c.Request(), c.Response())

	}()

	var actor *Node
	var box_id string
	var box_name string
	var box_type string
	defer func() {
		now := time.Now()
		diff := now.Sub(t1)

		go func(log_id string, c echo.Context, box_id string, box_name string, box_type string, connection_next string, diff time.Duration, order_box int, payload goja.Value) {
			if Config.DatabaseNflow.QueryInsertLog == "" {
				return
			}

			db, err := GetDB()
			if err != nil {
				log.Println(err)
				return
			}
			ctx := context.Background()
			conn, err := db.Conn(ctx)
			if err != nil {
				return
			}
			defer conn.Close()
			profile := GetProfile(c)
			username := ""
			if profile != nil {
				if _, ok := profile["username"]; ok {
					username = profile["username"]
				}
			}

			jsonPayload, err := json.Marshal(payload.Export())
			if err != nil {
				jsonPayload = []byte("{}")
			}
			ip := ""
			realip := ""
			url := ""
			userAgent := ""
			queryParam := ""
			hostname := ""
			host := ""

			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println(err)
					}
				}()
				ip = c.Request().RemoteAddr
				realip = c.RealIP()
				url = c.Request().URL.RawPath
				userAgent = c.Request().UserAgent()
				queryParam = c.Request().URL.Query().Encode()
				hostname = c.Request().URL.Hostname()
				host = c.Request().Host

			}()

			_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryInsertLog,
				log_id,                                  // $1
				box_id,                                  // $2
				box_name,                                // $3
				box_type,                                // $4
				url,                                     // $5
				username,                                // $6
				connection_next,                         // $7
				fmt.Sprintf("%dm", diff.Milliseconds()), // $8
				order_box,                               // $9
				string(jsonPayload),                     // $10
				ip,                                      // $11
				realip,                                  // $12
				userAgent,                               // $13
				queryParam,                              // $14
				hostname,                                // $15
				host,                                    // $16

			)
			if err != nil {
				log.Println(err)
			}

		}(log_id, c, box_id, box_name, box_type, connection_next, diff, order_box, payload)

		go func(c echo.Context, actor *Node, box_id string, box_name string, box_type string, connection_next string, diff time.Duration) {

			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()

			//log.Println(sbLog.String() + " - time:" + fmt.Sprint(diff))
			ctx := context.Background()
			db, err := GetDB()
			if err != nil {
				log.Println(err)
				return
			}
			conn, err := db.Conn(ctx)
			if err != nil {
				return
			}
			defer conn.Close()
			row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetApp, "app")
			var code string
			var json_code string
			err = row.Scan(&json_code, &code)
			if err != nil {
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

func (cc *Controller) Execute(c echo.Context, vm *goja.Runtime, next string, vars Vars, currentProcess *process.Process, payload goja.Value, fork bool) {
	var err error
	prev_box := ""
	if fork {
		fmt.Println("fork")
	}
	for next != "" {

		vm.Set("current_box", next)
		vm.Set("prev_box", prev_box)

		prev_box = next

		wg.Add(1)
		go func() {
			defer wg.Done()
			s, _ := session.Get("nflow_form", c)
			payload_map := make(map[string]interface{})

			if payload != nil {
				payload_map = payload.Export().(map[string]interface{})
			}

			for k, v := range s.Values {
				payload_map[k.(string)] = v
			}

			payload = vm.ToValue(payload_map)
		}()
		wg.Wait()

		next, payload, err = cc.step(c, vm, next, vars, currentProcess, payload)
		if err != nil {
			break
		}
		if fork {
			fmt.Println("fork")
		}

		// cut
		if payload != nil {
			if rawPayload, ok := payload.Export().(map[string]interface{}); ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s, _ := session.Get("nflow_form", c)
					for k, v := range rawPayload {
						s.Values[k] = v
					}

					s.Save(c.Request(), c.Response())
				}()
				wg.Wait()

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

	if next == "" && !fork {
		s, _ := session.Get("nflow_form", c)
		s.Values = make(map[interface{}]interface{})
		s.Save(c.Request(), c.Response())

		currentProcess.State = "end"
		currentProcess.Killeable = false
		currentProcess.Close()

	}

}
