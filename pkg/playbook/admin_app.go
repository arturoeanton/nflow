package playbook

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var FindNewApp map[string]bool = make(map[string]bool)

func SaveApp(c echo.Context) error {

	appJson := GetAppJsonFileName(c)
	pathBase := GetPathBase(c)
	FindNewApp[appJson] = true

	bytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Fatalln(err)
	}
	err1 := ioutil.WriteFile(pathBase+appJson+".json", bytes, 0644)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"msg": err1.Error()})
	}

	return c.JSON(200, echo.Map{"msg": "ok"})
}
