package playbook

import (
	"context"
	"log"

	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

func GetTemplateFromDB(param_name string) string {
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return ""
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Println(err)
		return ""
	}
	defer conn.Close()
	row := conn.QueryRowContext(context.Background(), Config.DatabaseNflow.QueryGetTemplate, param_name)

	var id int
	var name string
	var content string

	err = row.Scan(&id, &name, &content)
	if err != nil {
		log.Println(err)
		return ""
	}

	return content
}

func addFeatureTemplte(vm *goja.Runtime, c echo.Context) {

	vm.Set("get_template", func(param_name string) string {
		return GetTemplateFromDB(param_name)
	})

}
