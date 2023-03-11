package playbook

import (
	"encoding/json"
	"net/http"

	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

type StepPlugin struct {
}

func (s *StepPlugin) Run(cc *Controller, actor *Node, c echo.Context, vm *goja.Runtime, connection_next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error) {
	currentProcess.State = "run"
	currentProcess.Killeable = true
	name := actor.Data["dromedary_name"].(string)
	var payloadOut interface{}
	dataJs, _ := json.Marshal(actor.Data)
	payloadOut, next, err := Plugins[name].Run(c, vars, &payload, string(dataJs), nil)

	payload = vm.ToValue(payloadOut)
	currentProcess.Payload = payload
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
		currentProcess.State = "error"
		return "", payload, err
	}
	currentProcess.State = "end"
	if actor.Outputs != nil {
		if actor.Outputs[next] != nil {
			next = actor.Outputs[next].Connections[0].Node
		} else {
			connection_next = ""
		}
	}
	return next, payload, err
}
