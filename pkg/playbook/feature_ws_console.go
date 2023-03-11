package playbook

import (
	"log"

	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/dop251/goja"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var (
	WebsocketsConsole map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	WebsocketsAll     map[string]*websocket.Conn = make(map[string]*websocket.Conn)
)

func WebsocketConsole(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		uuidConn := uuid.New().String()
		WebsocketsConsole[uuidConn] = ws
		defer ws.Close()
		defer delete(WebsocketsConsole, uuidConn)
		for {
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				return
			}
			log.Println("cmd:" + msg)
			command(msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func SendConsoleWS(out string) {
	for _, ws := range WebsocketsConsole {
		err := websocket.Message.Send(ws, out)
		if err != nil {
			log.Println(err)
		}
	}
}

func SendWS(ws *websocket.Conn, out string) {
	err := websocket.Message.Send(ws, out)
	if err != nil {
		log.Println(err)
	}
}

func addFeatureWsConsole(vm *goja.Runtime, c echo.Context) {
	vm.Set("ws_console_log", func(data string) {
		SendConsoleWS(data)
	})
	vm.Set("ws_send", func(ws *websocket.Conn, data string) {
		SendWS(ws, data)
	})
}

func WSServer(c echo.Context, vm *goja.Runtime, next string, vars Vars, process *process.Process, cc *Controller) error {
	websocket.Handler(func(ws *websocket.Conn) {
		WebsocketsAll[process.UUID] = ws
		defer ws.Close()
		defer delete(WebsocketsAll, process.UUID)

		for {
			process.Type = "websocket"
			process.State = "recive"
			process.Ws = ws
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				return
			}
			payloadValue := map[string]interface{}{"msg": msg, "self": process.UUID, "websockets": WebsocketsAll}
			vm.Set("payload", payloadValue)
			payload := vm.Get("payload")
			cc.Execute(c, vm, next, vars, process, payload)

		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
