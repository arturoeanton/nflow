package plugins

import (
	"bytes"
	"html/template"
	"log"

	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
)

type TemplatePluings string

var (
	fxsTemplate map[string]interface{} = make(map[string]interface{})
)

func (d TemplatePluings) Run(c echo.Context,
	vars map[string]string, payload_in interface{}, dromedary_data string,
	callback chan string,
) (payload_out interface{}, next string, err error) {
	return nil, "output_1", nil
}

func init() {
	addFeatureTemplater()
}

func (d TemplatePluings) AddFeatureJS() map[string]interface{} {
	return fxsTemplate
}

func (d TemplatePluings) Name() string {
	return "template"
}

func templater(code string, data interface{}) string {
	t := template.Must(template.New("code").Parse(code))
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, data)
	return buf.String()
}

func mustacher(template string, data interface{}) string {

	ret, err := mustache.Render(template, data)
	if err != nil {
		log.Println(err)
		return ""
	}
	return ret
}

func addFeatureTemplater() {
	fxsTemplate["template"] = templater
	fxsTemplate["mustache"] = mustacher
}
