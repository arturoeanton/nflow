package commons

import (
	"log"

	customsession "github.com/arturoeanton/nFlow/pkg/nflow-session"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"

	"github.com/arturoeanton/nFlow/pkg/playbook"
)

func GetSessionStore(redisSessionConfig *playbook.RedisConfig) sessions.Store {
	if redisSessionConfig.Host != "" {
		// tls
		options_redis := []redis.DialOption{redis.DialUseTLS(redisSessionConfig.Tls), redis.DialTLSSkipVerify(redisSessionConfig.TlsSkipVerify)}
		store, err := customsession.NewRedisStore(redisSessionConfig.MaxConnectionPool, "tcp", redisSessionConfig.Host, redisSessionConfig.Password, options_redis) // set redis store
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
