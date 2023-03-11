package playbook

import (
	"github.com/labstack/echo/v4"
)

var PathBase string = "app/"

func GetPathBase(c echo.Context) string {
	base := "hopbox"
	if c.Param("app_name") != "" {
		base = c.Param("app_name")
	}
	return base + "/"
}

func GetAppJsonFileName(c echo.Context) string {
	base := "hopbox"
	if c.Param("app_name") != "" {
		base = c.Param("app_name")
	}
	return base
}
