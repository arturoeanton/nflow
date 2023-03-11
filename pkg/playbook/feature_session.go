package playbook

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func addFeatureSession(vm *goja.Runtime, c echo.Context) {
	vm.Set("set_session", func(name, k, v string) {
		s, _ := session.Get(name, c)
		s.Values[k] = v
		s.Save(c.Request(), c.Response())
	})

	vm.Set("get_session", func(name, k string) *string {
		s, _ := session.Get(name, c)
		r := fmt.Sprint(s.Values[k])
		return &r
	})

	vm.Set("open_session", func(name string) *map[string]interface{} {
		s, _ := session.Get(name, c)
		var r = make(map[string]interface{})
		for k, v := range s.Values {
			r[k.(string)] = v
		}
		return &r
	})

	vm.Set("save_session", func(name string, m map[string]interface{}) {
		s, _ := session.Get(name, c)
		for k, v := range m {
			s.Values[k] = v
		}
		s.Save(c.Request(), c.Response())
	})

}
