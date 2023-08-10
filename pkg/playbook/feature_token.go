package playbook

import (
	"context"
	"crypto/subtle"
	"log"

	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

func GetTokenFromDB(param_name string) []map[string]interface{} {
	db, err := GetDB()
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Println(err)
		return nil
	}
	defer conn.Close()
	rows, err := conn.QueryContext(context.Background(), Config.DatabaseNflow.QueryGetToken, param_name)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	var id int
	var name string
	var token string
	var start interface{}
	var expired interface{}
	var active bool
	var header string
	var token_type string
	ret := make([]map[string]interface{}, 0)
	found := false
	for rows.Next() {
		found = true
		err := rows.Scan(&id, &name, &token, &start, &expired, &active, &header, &token_type)
		if err != nil {
			log.Println(err)
			return nil
		}
		mapUser := make(map[string]interface{})
		mapUser["id"] = id
		mapUser["name"] = name
		mapUser["token"] = token
		mapUser["start"] = start
		mapUser["expired"] = expired
		mapUser["active"] = active
		mapUser["header"] = header
		mapUser["token_type"] = token_type
		ret = append(ret, mapUser)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	if !found {
		return nil
	}

	return ret
}

func ValidateTokenDB(name string, token string) bool {
	arrayToken := GetTokenFromDB(name)
	if arrayToken == nil {
		return false
	}
	for _, tokenMap := range arrayToken {
		if !tokenMap["active"].(bool) {
			return false
		}
		if subtle.ConstantTimeCompare([]byte(tokenMap["token_type"].(string)+" "+tokenMap["token"].(string)), []byte(token)) == 1 {
			return true
		}
	}
	return false
}

func addFeatureToken(vm *goja.Runtime, c echo.Context) {

	vm.Set("validate_token", func(name string, token string) bool {
		return ValidateTokenDB(name, token)
	})

	vm.Set("get_token", func(name string) []map[string]interface{} {
		ret := GetTokenFromDB(name)
		if ret == nil {
			return nil
		}
		return ret
	})

}
