package main

import (
	"crypto/subtle"
	"errors"
	"flag"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	"github.com/labstack/echo-contrib/session"

	"github.com/BurntSushi/toml"
	"github.com/arturoeanton/gocommons/utils"

	"github.com/arturoeanton/nFlow/pkg/commons"
	"github.com/arturoeanton/nFlow/pkg/playbook"
	"github.com/arturoeanton/nFlow/pkg/process"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var playbooks map[string]map[string]map[string]*playbook.Playbook = make(map[string]map[string]map[string]*playbook.Playbook)

func CheckError(c echo.Context, err error, code int) bool {
	if err != nil {
		c.JSON(code, echo.Map{
			"message": err.Error(),
			"code":    code,
		})
		return true
	}
	return false
}

func run(c echo.Context) error {
	ctx := c.Request().Context()
	db, err := playbook.GetDB()
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

	appJson := "app"
	endpoint := strings.Split(c.Request().RequestURI, "?")[0]

	v, ok := playbook.FindNewApp[appJson]
	if v || !ok {
		var err error
		playbooks[appJson], err = playbook.GetPlaybook(ctx, conn, appJson)
		if CheckError(c, err, 500) {
			return nil
		}
		playbook.FindNewApp[appJson] = false
	}

	nflow_next_node_run := ""
	if c.Request().Method == "POST" || c.Request().Method == "PUT" {
		if c.Request().FormValue("nflow_next_node_run") != "" {
			if c.Request().Form["nflow_next_node_run"] != nil {
				nflow_next_node_run = c.Request().Form["nflow_next_node_run"][0]
			}
		}
	}

	if nflow_next_node_run == "" {
		s, _ := session.Get("nflow_form", c)
		s.Values = make(map[interface{}]interface{})
		s.Save(c.Request(), c.Response())
	}

	runeable, vars, code, _, err := playbook.GetWorkflow(c, playbooks[appJson], endpoint, c.Request().Method, appJson)
	if CheckError(c, err, code) {
		return nil
	}

	uuid1 := uuid.New().String()
	e := runeable.Run(c, vars, nflow_next_node_run, uuid1, nil)
	return e
}

func main() {
	flag.Parse()
	configPath := "config.toml"
	if utils.Exists(configPath) {
		data, _ := utils.FileToString(configPath)
		if _, err := toml.Decode(data, &playbook.Config); err != nil {
			log.Println(err)
		}
	}

	playbook.LoadPlugins()

	playbook.RedisClient = redis.NewClient(&redis.Options{
		Addr:     playbook.Config.RedisConfig.Host,
		Password: playbook.Config.RedisConfig.Password, // no password set
		DB:       0,                                    // use default DB
	})

	e := echo.New()
	log.Println("URLBase:" + playbook.Config.URLConfig.URLBase)

	e.Use(middleware.Logger())
	e.Use(session.Middleware(commons.GetSessionStore(&playbook.Config.RedisSessionConfig)))

	e.Static("/site", "site/")
	e.File("/favicon.ico", "site/favicon.ico")
	e.File("/", "site/index.html")

	playbook.InitUI()

	e2 := echo.New()
	e2.Static("/site", "site/")
	e2.File("/favicon.ico", "site/favicon.ico")
	e2.File("/", "site/index.html")

	e2.Use(session.Middleware(commons.GetSessionStore(&playbook.Config.RedisSessionConfig)))
	gNFlow := e2.Group("/nflow")
	gNFlow.Use(session.Middleware(commons.GetSessionStore(&playbook.Config.RedisSessionConfig)))

	gNFlow.Static("/design", "design/")
	gNFlow.File("/favicon.ico", "design/favicon.ico")
	gNFlow.File("/", "design/index.html")
	gNFlow.GET("", playbook.Ui)
	gNFlow.GET("/", playbook.Ui)
	gNFlow.GET("/app", func(c echo.Context) error {
		ctx := c.Request().Context()
		db, err := playbook.GetDB()
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
		appJson := "app"
		rows, err := conn.QueryContext(ctx, playbook.Config.DatabaseNflow.QueryGetApp, appJson)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusNotFound, "Not Found")
			return nil
		}
		defer rows.Close()
		json := "{}"
		var default_js string
		for rows.Next() {
			err := rows.Scan(&json, &default_js)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusNotFound, "Not Found")
				return nil
			}
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
			c.HTML(http.StatusNotFound, "Not Found")
			return nil
		}

		c.HTML(200, json)
		return nil
	})
	gNFlow.POST("/app", playbook.SaveApp)

	gNFlow.Any("/modules", playbook.GetModules)

	gNFlow.Any("/ui/intellisense", playbook.Intellisense)

	gNFlow.GET("/module/manifest/:name", playbook.GetManifest)
	gNFlow.GET("/module/box/:name", playbook.GetBox)
	gNFlow.GET("/module/code/:name", playbook.GetCode)

	gNFlow.POST("/module/manifest/:name", playbook.PostManifest)
	gNFlow.POST("/module/box/:name", playbook.PostBox)
	gNFlow.POST("/module/code/:name", playbook.PostCode)

	gNFlow.DELETE("/module/:name", playbook.DeleteModule)

	//gNFlow.Any("/node/run/:flow_name/:node_id", runNode)

	gNFlow.Any("/process", process.GetProcesses)
	gNFlow.Any("/process/:wid", process.GetProcess)
	gNFlow.Any("/process/:wid/payload", process.GetProcessPayload)
	gNFlow.Any("/process/:wid/kill", process.KillWID)

	e2.Any("/*", run)

	e.Any("/*", run)

	if playbook.Config.HttpsDesingnerConfig.HTTPBasic {
		e2.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {

			user := playbook.GetUserFromDB(username)
			if user == nil {
				return false, errors.New("user not found")
			}
			validate := playbook.ValidateUserDB(username, password)
			if !validate {
				return false, errors.New("user not found")
			}
			if subtle.ConstantTimeCompare([]byte(user["rol"].(string)), []byte("ROL_DEV")) == 1 {
				return true, nil
			}

			return false, nil
		}))
	}

	if playbook.Config.HttpsEngineConfig.HTTPBasic {
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			user := playbook.GetUserFromDB(username)
			if user == nil {
				return false, errors.New("user not found")
			}
			validate := playbook.ValidateUserDB(username, password)
			if !validate {
				return false, errors.New("user not found")
			}
			return false, nil
		}))
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		RunServer(e2, &playbook.Config.HttpsDesingnerConfig)
	}()

	go func() {
		defer wg.Done()
		RunServer(e, &playbook.Config.HttpsEngineConfig)
	}()

	wg.Wait()

}

func RunServer(e *echo.Echo, httpsConfig *playbook.HttpsConfig) {

	e.HideBanner = true
	if httpsConfig.Enable {
		if httpsConfig.Cert == "" || httpsConfig.Key == "" {
			log.Println("(" + httpsConfig.Description + ")Starting server with auto TLS:" + httpsConfig.Address)
			e.StartAutoTLS(httpsConfig.Address)

		} else {
			log.Println("(" + httpsConfig.Description + ")Starting server with TLS:" + httpsConfig.Address)
			e.StartTLS(httpsConfig.Address, httpsConfig.Cert, httpsConfig.Key)
		}
	} else {
		log.Println("(" + httpsConfig.Description + ")Starting server without TLS:" + httpsConfig.Address)
		e.Start(httpsConfig.Address)
	}
}
