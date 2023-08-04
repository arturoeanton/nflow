package playbook

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

	jsonStr := ""

	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.QueryContext(ctx, Config.DatabaseNflow.QueryGetModules)
	if err != nil {
		return err
	}
	defer rows.Close()
	var form string
	var mod string
	var code sql.NullString
	var name string

	modules := make([]struct {
		Key    string `json:"key"`
		Module Module `json:"module"`
	}, 0) //make(map[string]Module)
	for rows.Next() {
		err := rows.Scan(&form, &mod, &code, &name)
		if err != nil {
			//continue
			log.Println(err)
			return err
		}

		module := Module{}
		jsonStr = mod
		err = json.Unmarshal([]byte(jsonStr), &module)
		if err != nil {
			flag := true
			module = Module{
				Title:    name,
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
			module.Title = name
		}
		if module.Editable == nil {
			flag := true
			module.Editable = &flag
		}
		if module.Custom {
			module.HTMLForm = form
		} else {
			module.HTMLForm = `<div>
		<div class="title-box"><i class="` + module.Icon + `"></i> ` + module.Title + `</div>
		<div class="box">
		  ` + form + `
		</div>
	  </div>`
		}
		//modules[name] = module
		modules = append(modules, struct {
			Key    string `json:"key"`
			Module Module `json:"module"`
		}{Key: name, Module: module})
	}

	c.JSON(http.StatusOK, modules)
	return nil
}

func GetManifest(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetModuleByName, name)
	var form string
	var mod string
	var code sql.NullString
	err = row.Scan(&form, &mod, &code)
	if err != nil {
		return err
	}
	c.String(http.StatusOK, mod)
	return nil
}

func GetBox(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetModuleByName, name)
	var form string
	var mod string
	var code sql.NullString
	err = row.Scan(&form, &mod, &code)
	if err != nil {
		return err
	}
	c.String(http.StatusOK, form)
	return nil
}

func GetCode(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	row := conn.QueryRowContext(ctx, Config.DatabaseNflow.QueryGetModuleByName, name)
	var form string
	var mod string
	var code sql.NullString
	err = row.Scan(&form, &mod, &code)
	if err != nil {
		return err
	}
	c.String(http.StatusOK, code.String)
	return nil
}

func PostManifest(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	selectCount := Config.DatabaseNflow.QueryCountModulesByName
	row := conn.QueryRowContext(ctx, selectCount, name)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	code := buf.String()
	if count == 0 {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryInsertModule, name, "", code, "")
	} else {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryUpdateModModuleByName, code, name)
	}
	if err != nil {
		return err
	}
	c.String(http.StatusOK, code)
	return nil
}

func PostBox(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	selectCount := Config.DatabaseNflow.QueryCountModulesByName
	row := conn.QueryRowContext(ctx, selectCount, name)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	code := buf.String()
	if count == 0 {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryInsertModule, name, code, "", "")
	} else {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryUpdateFormModuleByName, code, name)
	}
	if err != nil {
		return err
	}
	c.String(http.StatusOK, code)
	return nil
}

func PostCode(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	selectCount := Config.DatabaseNflow.QueryCountModulesByName
	row := conn.QueryRowContext(ctx, selectCount, name)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	code := buf.String()
	if count == 0 {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryInsertModule, name, "", "", code)
	} else {
		_, err = conn.ExecContext(ctx, Config.DatabaseNflow.QueryUpdateCodeModuleByName, code, name)
	}
	if err != nil {
		return err
	}
	c.String(http.StatusOK, code)
	return nil
}

func DeleteModule(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	name := c.Param("name")
	delete := Config.DatabaseNflow.QueryDeleteModule
	_, err = conn.ExecContext(ctx, delete, name)
	if err != nil {
		return err
	}
	c.NoContent(200)
	return nil
}
