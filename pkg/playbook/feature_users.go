package playbook

import (
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"log"

	"github.com/dop251/goja"
	"github.com/labstack/echo/v4"
)

func calcularSHA512_256(texto string) string {
	sha512_256 := sha512.New512_256()
	sha512_256.Write([]byte(texto))
	hashBytes := sha512_256.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func calcularSHA512_224(texto string) string {
	sha512_256 := sha512.New512_224()
	sha512_256.Write([]byte(texto))
	hashBytes := sha512_256.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func calcularSHA512(texto string) string {
	sha512 := sha512.New()
	sha512.Write([]byte(texto))
	hashBytes := sha512.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func GetUserFromDB(param_username string) map[string]interface{} {
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
	rows, err := conn.QueryContext(context.Background(), Config.DatabaseNflow.QueryGetUser, param_username)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	var id int
	var username string
	var password string
	var rol string
	var active bool
	found := false
	for rows.Next() {
		found = true
		err := rows.Scan(&id, &username, &password, &rol, &active)
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	if !found {
		return nil
	}

	mapUser := make(map[string]interface{})
	mapUser["id"] = id
	mapUser["username"] = username
	mapUser["password"] = password
	mapUser["rol"] = rol
	mapUser["active"] = active
	return mapUser
}

func ValidateUserDB(username string, password string) bool {
	user := GetUserFromDB(username)
	if user == nil {
		return false
	}
	if !user["active"].(bool) {
		return false
	}
	passwordSha512_256 := calcularSHA512_256(password)
	return subtle.ConstantTimeCompare([]byte(user["password"].(string)), []byte(passwordSha512_256)) == 1
}

func addFeatureUsers(vm *goja.Runtime, c echo.Context) {
	vm.Set("validate_user", func(username string, password string) bool {
		return ValidateUserDB(username, password)
	})
	vm.Set("get_user", func(username string) map[string]interface{} {
		return GetUserFromDB(username)
	})
	vm.Set("calcular_sha512_256", func(texto string) string {
		return calcularSHA512_256(texto)
	})
	vm.Set("calcular_sha512_224", func(texto string) string {
		return calcularSHA512_224(texto)
	})
	vm.Set("calcular_sha512", func(texto string) string {
		return calcularSHA512(texto)
	})

}
