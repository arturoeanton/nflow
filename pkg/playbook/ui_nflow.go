package playbook

import (
	"fmt"
	"net/http"

	"github.com/arturoeanton/gocommons/utils"
	"github.com/labstack/echo/v4"
)

func Ui(c echo.Context) error {
	//fmt.Println("header", c.Request().Header)
	appName := c.Param("app_name")
	if appName == ":app_name" {
		appName = "app"
		c.Redirect(http.StatusTemporaryRedirect, "/app/nflow")
	}
	ui, _ := utils.FileToString("design/index.html")
	c.HTML(200, ui)
	return nil
}

func InitUI() {
	jsVars["ws_console_log()"] = "log in nflow console"
	jsVars["ws_send()"] = "send by websocket"
	jsVars["path_base"] = "Directory root"
	jsVars["config"] = "Config of nflow"
	jsVars["ctx"] = "context.Background()"
	jsVars["url_base"] = "url base"
	jsVars["wkill()"] = "kill workflow"

	jsVars["wid"] = "Instance id"

	jsVars["vars"] = "Path variable"
	jsVars["path_vars"] = "Path variable"

	jsVars["c"] = "Context echo"
	jsVars["echo_context"] = "Context echo"
	jsVars["payload"] = "payload"
	jsVars["next"] = "next edge example (output_1, output_2)"
	jsVars["dromedary_data"] = "param of box"
	jsVars["nflow_data"] = "param of box"
	jsVars["__outputs"] = ""
	jsVars["__flow_name"] = "name of  flow (playbook)"
	jsVars["current_box"] = "current box id"
	jsVars["prev_box"] = "previous box id"

}

func Intellisense(c echo.Context) error {

	jsWords := make([]string, 0)
	for _, p := range Plugins {
		for key := range p.AddFeatureJS() {
			jsWords = append(jsWords, fmt.Sprintf("%s()", key))
		}
	}

	for key := range jsVars {
		jsWords = append(jsWords, fmt.Sprintf("%s", key))
	}

	jsWords = append(jsWords, "function (){\n}")
	jsWords = append(jsWords, "function main(){\n}")
	jsWords = append(jsWords, "function ")
	jsWords = append(jsWords, "if ()")
	jsWords = append(jsWords, "else ")
	jsWords = append(jsWords, "for (var i in $list){}")
	jsWords = append(jsWords, "for ")
	jsWords = append(jsWords, "while ")
	jsWords = append(jsWords, "var ")

	jsWords = append(jsWords, "c.JSON(200,{})")
	jsWords = append(jsWords, "JSON.parse")
	jsWords = append(jsWords, "JSON.stringify")

	return c.JSON(200, map[string][]string{"js_words": jsWords})
}
