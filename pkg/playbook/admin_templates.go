package playbook

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeleteTemplateByName(c echo.Context) error {
	name := c.Param("name")
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	ctx := c.Request().Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryDeleteTemplate, name)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK"})
}
func GetTemplateByName(c echo.Context) error {
	name := c.Param("name")
	template := GetTemplateFromDB(name)

	content := template
	return c.JSON(http.StatusOK, map[string]interface{}{"name": name, "content": content})
}

func GetAllTemplates(c echo.Context) error {
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})

	}
	ctx := c.Request().Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})

	}
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, Config.DatabaseNflow.QueryGetTemplates)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})

	}
	defer rows.Close()

	ret := make([]map[string]interface{}, 0)
	for rows.Next() {
		var name string
		var content string
		var id int
		err = rows.Scan(&id, &name, &content)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})

		}
		ret = append(ret, map[string]interface{}{"id": id, "name": name, "content": content})
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	return c.JSON(http.StatusOK, ret)
}

func UpdateTemplate(c echo.Context) error {
	jsonValue := make(map[string]interface{})
	err := c.Bind(&jsonValue)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	name := jsonValue["name"].(string)
	content := jsonValue["content"].(string)
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	ctx := c.Request().Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetTemplateCount, name)
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	if count > 0 {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryUpdateTemplate, content, name)
	} else {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryInsertTemplate, content, name)
	}
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK"})
}

func CreateTemplate(c echo.Context) error {
	jsonValue := make(map[string]interface{})
	err := c.Bind(&jsonValue)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	name := jsonValue["name"].(string)
	content := jsonValue["content"].(string)
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	ctx := c.Request().Context()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetTemplateCount, name)
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	if count == 0 {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryInsertTemplate, content, name)
	} else {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]interface{}{"message": "Not Found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OK"})
}
