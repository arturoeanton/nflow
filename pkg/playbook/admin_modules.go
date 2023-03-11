package playbook

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/arturoeanton/gocommons/utils"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Title     string                 `json:"title"`
	Icon      string                 `json:"icon"`
	In        int                    `json:"in"`
	Out       int                    `json:"out"`
	Editable  *bool                  `json:"editable"`
	Hide      bool                   `json:"hide"`
	Custom    bool                   `json:"custom"`
	Param     map[string]interface{} `json:"param"`
	BoxColor  string                 `json:"boxcolor"`
	FontColor string                 `json:"fontcolor"`
	HTMLForm  string
}

func GetModules(c echo.Context) error {
	pathBase := GetPathBase(c)
	modules := make(map[string]Module)

	jsonStr := ""
	files, err := ioutil.ReadDir(pathBase + "modules/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	for _, f := range files {

		if !f.IsDir() {
			continue
		}

		module := Module{}
		jsonStr, _ = utils.FileToString(pathBase + "modules/" + f.Name() + "/mod.json")
		err = json.Unmarshal([]byte(jsonStr), &module)
		if err != nil {
			flag := true
			module = Module{
				Title:    f.Name(),
				Icon:     "fas fa-bomb",
				In:       0,
				Out:      0,
				Editable: &flag,
				HTMLForm: "",
			}
		}
		if module.Hide {
			continue
		}
		if module.Title == "" {
			module.Title = f.Name()
		}
		if module.Editable == nil {
			flag := true
			module.Editable = &flag
		}
		if module.Custom {
			module.HTMLForm, _ = utils.FileToString(pathBase + "modules/" + f.Name() + "/" + f.Name() + ".html")
		} else {
			html, _ := utils.FileToString(pathBase + "modules/" + f.Name() + "/" + f.Name() + ".html")
			module.HTMLForm = `<div>
		<div class="title-box"><i class="` + module.Icon + `"></i> ` + module.Title + `</div>
		<div class="box">
		  ` + html + `
		</div>
	  </div>`
		}
		modules[f.Name()] = module
	}

	c.JSON(http.StatusOK, modules)
	return nil
}

func GetManifest(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	code, _ := utils.FileToString(pathBase + "modules/" + name + "/mod.json")
	c.String(http.StatusOK, code)
	return nil
}

func GetBox(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	code, _ := utils.FileToString(pathBase + "modules/" + name + "/" + name + ".html")
	c.String(http.StatusOK, code)
	return nil
}

func GetCode(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	code, _ := utils.FileToString(pathBase + "modules/" + name + "/js/code.js")
	c.String(http.StatusOK, code)
	return nil
}

func PostManifest(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	path := pathBase + "modules/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	code := buf.String()
	utils.StringToFile(path+"/mod.json", code)
	c.String(http.StatusOK, code)
	return nil
}

func PostBox(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	path := pathBase + "modules/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	code := buf.String()
	utils.StringToFile(path+"/"+name+".html", code)
	c.String(http.StatusOK, code)
	return nil
}

func PostCode(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	path := pathBase + "modules/" + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	path = path + "/js"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	code := buf.String()
	utils.StringToFile(path+"/code.js", code)
	c.String(http.StatusOK, code)
	return nil
}

func DeleteModule(c echo.Context) error {
	pathBase := GetPathBase(c)
	name := c.Param("name")
	path := pathBase + "modules/" + name
	os.RemoveAll(path)
	c.NoContent(200)
	return nil
}
