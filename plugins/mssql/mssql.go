package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

type dromedary string

type ConfigMSSQL struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Port     string `toml:"port"`
	Server   string `toml:"server"`
}

var fxs map[string]interface{} = make(map[string]interface{})

func (d dromedary) Run(c echo.Context,
	vars map[string]string, payload_in interface{}, dromedary_data string,
	callback chan string,
) (payload_out interface{}, next string, err error) {
	return nil, "output_1", nil
}

func init() {
	addFeatureCommon()
	log.Println("Started mssql")
}

func (d dromedary) AddFeatureJS() map[string]interface{} {
	return fxs
}

func (d dromedary) Name() string {
	return "mssql"
}

var Dromedary dromedary

//go build -buildmode=plugin -o demo2.so
