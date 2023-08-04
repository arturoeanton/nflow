package playbook

import (
	"bufio"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var FindNewApp map[string]bool = make(map[string]bool)

func SaveApp(c echo.Context) error {
	appJson := "app"
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "Not Found")
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusNotFound, "Not Found")
		return nil
	}
	defer conn.Close()

	scanner := bufio.NewScanner(c.Request().Body)

	// Creamos un Builder para concatenar las líneas en un string
	var builder strings.Builder

	// Iteramos sobre cada línea y la agregamos al Builder
	for scanner.Scan() {
		builder.WriteString(scanner.Text())
	}

	// Comprobamos si hubo algún error en el scanner
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Obtenemos el string resultante
	result := builder.String()

	if err != nil {
		c.JSON(http.StatusNotFound, echo.Map{"msg": err.Error()})
		return nil
	}

	resul, err := conn.ExecContext(ctx, Config.DatabaseNflow.QueryUpdateApp, result, appJson)

	if err != nil {
		conn2, err := db.Conn(ctx)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusNotFound, "Not Found")
			return nil
		}
		defer conn2.Close()
		resul, err = conn2.ExecContext(ctx, Config.DatabaseNflow.QueryUpdateApp, result, appJson)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, echo.Map{"msg": err.Error()})
			return nil
		}
	}
	_, err = resul.RowsAffected()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, echo.Map{"msg": err.Error()})
		return nil
	}

	FindNewApp[appJson] = true

	return c.JSON(200, echo.Map{"msg": "ok"})
}
