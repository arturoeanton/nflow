package main

import (
	"crypto/subtle"
	"flag"
	"github.com/go-redis/redis"
	"github.com/labstack/echo-contrib/session"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/arturoeanton/gocommons/utils"
	"github.com/gorilla/sessions"
	customsession "github.com/piggyman007/echo-session"

	"github.com/arturoeanton/nFlow/pkg/playbook"
	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

/*
func runNode(c echo.Context) error {
	pathBase := playbook.GetPathBase(c)
	appJson := playbook.GetAppJsonFileName(c)
	v, ok := playbook.FindNewApp[appJson]
	if v || !ok {
		var err error
		playbooks[appJson], err = playbook.GetPlaybook(pathBase, appJson)
		if CheckError(c, err, 500) {
			return nil
		}
		playbook.FindNewApp[appJson] = false
	}

	cc := playbook.Controller{
		Methods:  []string{c.Request().Method},
		Playbook: playbooks[appJson][c.Param("flow_name")]["data"],
		FlowName: c.Param("flow_name"),
		AppName:  appJson,
	}
	vars := playbook.Vars{}
	for k := range c.QueryParams() {
		vars[k] = c.QueryParam(k)
	}

	uuid1 := uuid.New().String()
	return cc.Run(c, vars, c.Param("node_id"), uuid1, nil)
}*/

func run(c echo.Context) error {
	pathBase := playbook.GetPathBase(c)
	appJson := playbook.GetAppJsonFileName(c)
	endpoint := strings.Split(c.Request().RequestURI, "?")[0]

	v, ok := playbook.FindNewApp[appJson]
	if v || !ok {
		var err error
		playbooks[appJson], err = playbook.GetPlaybook(pathBase, appJson)
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

	runeable, vars, err, code := playbook.GetWorkflow(c, playbooks[appJson], endpoint, c.Request().Method, appJson)
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
	e.Use(session.Middleware(GetSessionStore(&playbook.Config.RedisSessionConfig)))

	e.Static("/site", "site/")
	e.File("/favicon.ico", "site/favicon.ico")
	e.File("/", "site/index.html")

	playbook.InitUI()

	e2 := echo.New()
	e2.Static("/site", "site/")
	e2.File("/favicon.ico", "site/favicon.ico")
	e2.File("/", "site/index.html")

	e2.Use(session.Middleware(GetSessionStore(&playbook.Config.RedisSessionConfig)))
	gNFlow := e2.Group("/nflow")
	gNFlow.Use(session.Middleware(GetSessionStore(&playbook.Config.RedisSessionConfig)))

	gNFlow.Static("/design", "design/")
	gNFlow.File("/favicon.ico", "design/favicon.ico")
	gNFlow.File("/", "design/index.html")
	gNFlow.GET("", playbook.Ui)
	gNFlow.GET("/", playbook.Ui)
	gNFlow.GET("/app", func(c echo.Context) error {
		pathBase := playbook.GetPathBase(c)
		appJson := playbook.GetAppJsonFileName(c)

		if !utils.Exists(pathBase + appJson + ".json") {
			c.HTML(http.StatusNotFound, "Not Found")
			return nil
			/*copy.Copy("app_template", pathBase)
			content := `
			{
				"drawflow": {
					"Home": {
						"data": {}
					},
					"": {
						"data": {}
					}
				}
			}
			`
			utils.StringToFile(pathBase+appJson+".json", content)*/
		}
		content, _ := utils.FileToString(pathBase + appJson + ".json")
		c.HTML(200, content)
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

	if playbook.Config.HttpsDesingnerConfig.User != nil {
		e2.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte(*playbook.Config.HttpsDesingnerConfig.User)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(*playbook.Config.HttpsDesingnerConfig.Password)) == 1 {
				return true, nil
			}
			return false, nil
		}))
	}

	if playbook.Config.HttpsEngineConfig.User != nil {
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte(*playbook.Config.HttpsEngineConfig.User)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(*playbook.Config.HttpsEngineConfig.Password)) == 1 {
				return true, nil
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

func GetSessionStore(redisSessionConfig *playbook.RedisConfig) sessions.Store {
	if redisSessionConfig.Host != "" {
		store, err := customsession.NewRedisStore(redisSessionConfig.MaxConnectionPool, "tcp", redisSessionConfig.Host, redisSessionConfig.Password) // set redis store
		if err != nil {
			log.Printf("could not create redis store: %s - using cookie store instead", err.Error())
			return sessions.NewCookieStore([]byte("secret"))
		}
		opts := customsession.Options{
			MaxAge:   3600, // session timeout in seconds
			Secure:   true, // secure cookie flag
			HttpOnly: true, // httponly flag
		}
		store.Options(opts)
		return store
	}
	return sessions.NewCookieStore([]byte("secret"))
}
