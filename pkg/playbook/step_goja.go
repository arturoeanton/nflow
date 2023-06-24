package playbook

import (
	"log"
	"net/http"

	"github.com/arturoeanton/gocommons/utils"
	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/google/uuid"
	babel "github.com/jvatic/goja-babel"
	"github.com/labstack/echo/v4"
)

var (
	semVM     = make(chan int, 50) // 70 aguanto - 80 no aguanto
	isem  int = 0
)

type StepJS struct {
}

func (s *StepJS) Run(cc *Controller, actor *Node, c echo.Context, vm *goja.Runtime, connection_next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error) {
	pathBase := GetPathBase(c)

	currentProcess.State = "run"
	currentProcess.Killeable = true
	code := "function main(){}"
	actor.Data["storage_id"] = uuid.New().String()
	if _, ok := actor.Data["compile"]; !ok {

		if _, ok := actor.Data["script"]; ok {
			filename := pathBase + actor.Data["script"].(string) + ".js"
			if utils.Exists(filename) {
				code, _ = utils.FileToString(filename)
				code = babelTransform(code)
				actor.Data["compile"] = code
			}
		}
		if _, ok := actor.Data["code"]; ok {
			code = actor.Data["code"].(string)
			code = babelTransform(code)
			actor.Data["compile"] = code
		}
	}

	if _, ok := actor.Data["compile"]; ok {
		code = actor.Data["compile"].(string)
	}
	code = code + "\nmain()"

	outputs := make(map[string]string)
	for key, o := range actor.Outputs {
		outputs[key] = o.Connections[0].Node
	}

	vm.Set("payload", payload)
	if payload == nil || payload.Equals(goja.NaN()) || payload.Equals(goja.Null()) || goja.IsUndefined(payload) {
		vm.Set("payload", make(map[string]interface{}))
	}

	vm.Set("next", connection_next)
	vm.Set("dromedary_data", actor.Data)
	vm.Set("nflow_data", actor.Data)
	vm.Set("__outputs", outputs)
	vm.Set("__flow_name", cc.FlowName)
	vm.Set("__flow_app", cc.AppName)

	err := func() error {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("runJs_00010 ****", err)
			}
		}()
		semVM <- 1
		isem++
		_, err := vm.RunString(code)
		isem--
		<-semVM
		return err
	}()

	if err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, echo.Map{
				"message": err.Error(),
				"actor":   actor,
			})
			currentProcess.State = "error"
			return "", payload, err
		}
	}
	payload = vm.Get("payload")
	currentProcess.Payload = payload.Export()
	connection_next = vm.Get("next").String()
	currentProcess.State = "end"
	if actor.Outputs != nil {
		if actor.Outputs[connection_next] != nil {
			connection_next = actor.Outputs[connection_next].Connections[0].Node
		} else {
			connection_next = ""
		}
	}
	return connection_next, payload, nil
}

func babelTransform(code string) string {
	/*t1 := time.Now()
	defer func() {
		diff := time.Now().Sub(t1)
		log.Println("babel time:" + fmt.Sprint(diff))
	}()*/
	babel.Init(4) // Setup 4 transformers (can be any number > 0)
	res, err := babel.TransformString(
		code,
		map[string]interface{}{
			"plugins": []string{
				"transform-block-scoping",
				"transform-block-scoped-functions",
				"transform-arrow-functions",
				"transform-classes",
				"transform-computed-properties",
				"transform-destructuring",
				"transform-for-of",
				"transform-template-literals",
				"transform-parameters",
				"transform-spread",
				"transform-shorthand-properties",
				"transform-duplicate-keys",
				"transform-object-super",
				"transform-literals",
				"transform-function-name",
				"transform-sticky-regex",
				"transform-typeof-symbol",
				"transform-unicode-regex",
			},
		},
	)
	if err != nil {
		log.Println(code)
		log.Println(err)
		log.Println("babelTransform_00010 ****", err)
	}
	return res
}
