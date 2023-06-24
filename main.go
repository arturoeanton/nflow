package main

import (
	"flag"
	"log"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"

	"github.com/arturoeanton/gocommons/utils"

	"github.com/arturoeanton/nFlow/pkg/playbook"
	"github.com/arturoeanton/nFlow/pkg/process"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-redis/redis"
	"github.com/otiai10/copy"
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
}

func run(c echo.Context) error {
	pathBase := playbook.GetPathBase(c)
	appJson := playbook.GetAppJsonFileName(c)
	endpoint := strings.Split(c.Request().RequestURI, "?")[0][len(appJson):]

	v, ok := playbook.FindNewApp[appJson]
	if v || !ok {
		var err error
		playbooks[appJson], err = playbook.GetPlaybook(pathBase, appJson)
		if CheckError(c, err, 500) {
			return nil
		}
		playbook.FindNewApp[appJson] = false
	}
	runeable, vars, err, code := playbook.GetWorkflow(playbooks[appJson], endpoint, c.Request().Method, appJson)
	if CheckError(c, err, code) {
		return nil
	}

	uuid1 := uuid.New().String()
	e := runeable.Run(c, vars, "", uuid1, nil)

	return e
}

func main() {
	workspace := flag.String("w", "app/", "folder of workspace")
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
	log.Println("PathBase (workspace):" + *workspace)
	log.Println("URLBase:" + playbook.Config.URLConfig.URLBase)

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("/site", *workspace+"site/")
	e.File("/favicon.ico", *workspace+"site/favicon.ico")
	e.File("/", *workspace+"site/index.html")

	e.GET("/__health", func(c echo.Context) error { return c.JSON(200, echo.Map{"alive": true}) })

	playbook.InitUI()

	e2 := echo.New()
	e2.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	gNFlow := e2.Group("/:app_name/nflow")
	gNFlow.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	gNFlow.Static("/site", "site/")
	gNFlow.File("/favicon.ico", "site/favicon.ico")
	gNFlow.File("/", "site/index.html")
	gNFlow.GET("", playbook.Ui)
	gNFlow.GET("/", playbook.Ui)
	gNFlow.GET("/app", func(c echo.Context) error {
		pathBase := playbook.GetPathBase(c)
		appJson := playbook.GetAppJsonFileName(c)

		if !utils.Exists(pathBase + appJson + ".json") {
			copy.Copy("app_template", pathBase)
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
			utils.StringToFile(pathBase+appJson+".json", content)
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

	gNFlow.Any("/node/run/:flow_name/:node_id", runNode)

	gNFlow.GET("/console/ws", playbook.WebsocketConsole)

	gNFlow.Any("/process", process.GetProcesses)
	gNFlow.Any("/process/:wid", process.GetProcess)
	gNFlow.Any("/process/:wid/payload", process.GetProcessPayload)
	gNFlow.Any("/process/:wid/kill", process.KillWID)

	e2.Any("/:app_name/*", run)

	e.Any("/:app_name/*", run)

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
