package playbook

import (
	"github.com/labstack/echo/v4"
)

var PathBase string = "app/"

func GetPathBase(c echo.Context) string {
	base := "app"

	return base + "/"
}

func GetAppJsonFileName(c echo.Context) string {
	base := "app"

	return base
}
