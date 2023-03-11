package playbook

import (
	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

type StepWebsocket struct {
}

func (s *StepWebsocket) Run(cc *Controller, actor *Node, c echo.Context, vm *goja.Runtime, connection_next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error) {
	if actor.Outputs != nil {
		if actor.Outputs[connection_next] != nil {
			currentProcess.State = "start"
			currentProcess.Killeable = true
			nextWS := actor.Outputs[connection_next].Connections[0].Node
			WSServer(c, vm, nextWS, vars, currentProcess, cc)
		}
	}
	return "", payload, nil
}
