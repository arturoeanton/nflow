package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	Config     ConfigMSSQL
	ConnString string
)

type any = interface{}

func ConfigDatabase(user string, password string, database string, server string, port string) string {
	Config.Database = database
	Config.User = user
	Config.Password = password
	Config.Server = server
	Config.Port = port
	ConnString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)
	return ConnString
}

func Exce(tsql string, args ...interface{}) interface{} {
	type Result struct {
		RowsAffected int64
		LastInsertId int64
		Error        string
	}
	db, err := sql.Open("mssql", ConnString)
	if err != nil {
		errStr := fmt.Sprintf("Open connection failed: %s", err.Error())
		log.Println(errStr)
		return Result{
			RowsAffected: -1, LastInsertId: -1, Error: errStr,
		}
	}
	defer db.Close()

	result, err1 := db.Exec(tsql, args...)
	if err1 != nil {
		fmt.Println("Error operation row: " + err1.Error())
		return Result{
			RowsAffected: -1, LastInsertId: -1, Error: err1.Error(),
		}
	}

	rowsAffected, _ := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId()
	return Result{
		RowsAffected: rowsAffected, LastInsertId: lastInsertId, Error: "ok",
	}
}

func Query(tsql string, args ...interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	db, err := sql.Open("mssql", ConnString)
	if err != nil {
		errStr := fmt.Sprintf("Open connection failed: %s", err.Error())
		log.Println(errStr)
		result["error"] = errStr
		return result
	}
	defer db.Close()

	rows, err1 := db.Query(tsql, args...)
	if err1 != nil {
		fmt.Println("Error reading rows: " + err1.Error())
		result["error"] = err1.Error()
		return result
	}
	defer rows.Close()
	var count int64 = 0
	colums, err2 := rows.Columns()
	if err2 != nil {
		fmt.Println("Get colums: " + err2.Error())
		result["error"] = err2.Error()
		return result
	}
	resultRows := make([][]interface{}, 0)
	lencolumns := len(colums)

	for rows.Next() {
		row := make([]interface{}, lencolumns)
		scanArgs := make([]interface{}, lencolumns)
		for i := range row {
			scanArgs[i] = &row[i]
		}
		fmt.Println(count)
		count++

		e := rows.Scan(scanArgs...)
		fmt.Println(e)

		resultRows = append(resultRows, row)
		fmt.Println(row)
	}
	result["columns"] = colums
	result["rows"] = resultRows
	result["count"] = count
	return result
}

func addFeatureCommon() {

	fxs["mssql_exec"] = Exce
	fxs["mssql_query"] = Query
	fxs["mssql_config"] = ConfigDatabase

	//fxs["connect"] = Connect
	//fxs["disconnect"] = Disconnect

}

func main() {
	str := ConfigDatabase("sa", "MyPass@word", "test", "localhost", "1433")
	Exce("insert into dbo.t1 values (?, ?, ?)", "dato3", "dato3", "3")
	r := Query("select * from dbo.t1")
	fmt.Println(str)
	fmt.Println(r)
}
