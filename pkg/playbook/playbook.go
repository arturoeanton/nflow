package playbook

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (

	// Config is ...
	Config ConfigWorkspace
	// PathBase is ...
)

func GetPlaybook(ctx context.Context, conn *sql.Conn, pbName string) (map[string]map[string]*Playbook, error) {
	rows, err := conn.QueryContext(ctx, Config.DatabaseNflow.QueryGetApp, pbName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var flow_json string
	var default_js string
	for rows.Next() {
		err := rows.Scan(&flow_json, &default_js)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	data := make(map[string]map[string]map[string]*Playbook)
	err = json.Unmarshal([]byte(flow_json), &data)
	if err != nil {
		return nil, err
	}

	return data["drawflow"], nil
}

func comparePath(template string, real string) (bool, Vars) {
	termsOfTemplate := strings.Split(template, "/")
	termsOfReal := strings.Split(real, "/")
	vars := make(Vars)
	if len(termsOfTemplate) != len(termsOfReal) {
		return false, nil
	}
	for i, tt := range termsOfTemplate {
		if tt == "" {
			continue
		}
		tr := termsOfReal[i]
		if tt[0] == ':' {
			vars[tt[1:]] = tr
			continue
		}
		if tt != tr {
			return false, nil
		}
	}
	return true, vars
}

func GetWorkflow(c echo.Context, playbooks map[string]map[string]*Playbook, wfPath string, method string, appName string) (Runeable, Vars, error, int) {

	for key, flows := range playbooks {
		for _, pb := range flows {
			for _, item := range *pb {
				data := item.Data
				typeItem := data["type"]

				if typeItem == "starter" {

					methodItem := data["method"]
					if methodItem != "ANY" {
						if methodItem != method {
							continue
						}
					}
					urlpattern := data["urlpattern"].(string)
					flag, vars := comparePath(urlpattern, wfPath)
					if flag {
						c := Controller{
							Methods:  []string{method},
							Start:    item,
							Playbook: pb,
							FlowName: key,
							AppName:  appName,
						}

						return Runeable(&c), vars, nil, http.StatusOK
					}
				}
			}
		}
	}

	/*
		for key, c := range pb.Controllers {
			flag, vars := comparePath(key, wfPath)

			if !common.ContainsString(c.GetMethods(), method) {
				continue
			}

			if flag {
				return Runeable(&c), vars, nil, http.StatusOK
			}
		}
	*/
	return nil, nil, errors.New("not found"), http.StatusNotFound
}
