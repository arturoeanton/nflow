package playbook

import (
	"context"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/go-redis/redis"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	RedisClient *redis.Client
)

func addGlobals(vm *goja.Runtime, c echo.Context) {
	pathBase := GetPathBase(c)
	form, err1 := c.FormParams()
	if err1 != nil {
		vm.Set("form", make(map[string][]string))
	} else {
		vm.Set("form", (map[string][]string)(form))
	}

	s, _ := session.Get("nflow_form", c)
	for k, v := range form {
		if len(v) == 1 {
			s.Values[k] = v[0]
			continue
		}
		s.Values[k] = v
	}
	s.Save(c.Request(), c.Response())

	vm.Set("profile", nil)
	//sess_auth, _ := app.Store.Get(c.Request(), "auth-session")
	//if val, ok := sess_auth.Values["profile"]; ok {
	//	vm.Set("profile", val)
	//}

	vm.Set("redis_hset", RedisClient.HSet)
	vm.Set("redis_hget", RedisClient.HGet)
	vm.Set("redis_hdel", RedisClient.HDel)
	vm.Set("redis_expire", func(key string, s int32) {
		RedisClient.Expire(key, time.Duration(s)*time.Second)
	})

	fmt.Println("REDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDISREDIS")
	vm.Set("path_base", pathBase)
	vm.Set("config", Config)
	vm.Set("env", Config.Env)

	vm.Set("url_base", Config.URLConfig.URLBase)
	vm.Set("__vm", *vm)
	vm.Set("ctx", context.Background())

}
