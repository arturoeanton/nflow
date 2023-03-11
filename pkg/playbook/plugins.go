package playbook

import (
	"fmt"
	"plugin"

	"github.com/labstack/echo/v4"

	"github.com/arturoeanton/nFlow/pkg/plugins"
)

type NflowPlugin interface {
	Run(c echo.Context, vars map[string]string, payload_in interface{}, dromedary_data string, callback chan string) (payload_out interface{}, next string, err error)
	Name() string
	AddFeatureJS() map[string]interface{}
}

var Plugins map[string]NflowPlugin

func LoadPlugins() {
	Plugins = make(map[string]NflowPlugin)

	pluing1 := plugins.ClientHTTP("client_http")
	Plugins[pluing1.Name()] = pluing1

	pluing2 := plugins.GojaPlugin("goja")
	Plugins[pluing2.Name()] = pluing2

	pluing3 := plugins.TemplatePluings("template")
	Plugins[pluing3.Name()] = pluing3

	pluing4 := plugins.MailPlugin("mail")
	Plugins[pluing4.Name()] = pluing4

	pluing5 := plugins.RulePlugin("rule")
	Plugins[pluing5.Name()] = pluing5

	for _, mod := range Config.PluginConfig.Plugins {

		// load module
		// 1. open the so file to load the symbols
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Try use \n\tmake run")
					fmt.Println("Recovered in f", r)
				}
			}()
			plug, err := plugin.Open(mod)
			if err != nil {
				fmt.Println(err)
				//os.Exit(1)
			}

			// 2. look up a symbol (an exported function or variable)
			// in this case, variable NflowPlugin
			symNflowPlugin, err := plug.Lookup("Dromedary")
			if err != nil {
				fmt.Println(err)
				//os.Exit(1)
			}

			// 3. Assert that loaded symbol is of a desired type
			// in this case interface type Dromedary (defined above)
			var plugin NflowPlugin
			plugin, ok := symNflowPlugin.(NflowPlugin)
			if !ok {
				fmt.Println("unexpected type from module symbol")
				//os.Exit(1)
			}
			Plugins[plugin.Name()] = plugin
			fmt.Println(plugin.Name())
		}()
	}

}
