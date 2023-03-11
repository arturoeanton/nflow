package playbook

import (
	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

var (
	Steps map[string]Step = make(map[string]Step)
)

type Step interface {
	Run(cc *Controller, actor *Node, c echo.Context, vm *goja.Runtime, connection_next string, vars Vars, currentProcess *process.Process, payload goja.Value) (string, goja.Value, error)
}

func init() {
	Steps["gorutine"] = &StepGorutine{}
	Steps["js"] = &StepJS{}
	Steps["dromedary"] = &StepPlugin{}
	Steps["dromedary_callback"] = &StepPluginCallback{}
	Steps["websocket"] = &StepWebsocket{}
}
