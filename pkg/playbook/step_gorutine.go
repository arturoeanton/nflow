package playbook

import (
	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type StepGorutine struct {
}

func (s *StepGorutine) Run(cc *Controller, actor *Node, c echo.Context, vm *goja.Runtime, connection_next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error) {
	currentProcess.State = "run"
	if actor.Outputs["output_2"] != nil {
		next2 := actor.Outputs["output_2"].Connections[0].Node
		uuid2 := uuid.New().String()
		c.Response().Header().Add("Dromedary-Wid-2", uuid2)
		go cc.RunWithCallback(c, vars, next2, uuid2, nil)
	}
	connection_next = actor.Outputs[connection_next].Connections[0].Node
	currentProcess.State = "end"
	return connection_next, payload, nil
}
