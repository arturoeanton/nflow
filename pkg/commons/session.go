package commons

import (
	"log"

	"github.com/gorilla/sessions"
	customsession "github.com/piggyman007/echo-session"

	"github.com/arturoeanton/nFlow/pkg/playbook"
)

func GetSessionStore(pgSessionConfig *playbook.PgSessionConfig) sessions.Store {

	if pgSessionConfig.Url != "" {
		log.Println("pg session")
		store, err := customsession.NewPostgresStore(
			pgSessionConfig.Url,
			[]byte("secret"),
		)
		if err != nil {
			log.Println(err)
			return nil
		}
		return store
	}
	/*
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

	*/
	return sessions.NewCookieStore([]byte("secret"))
}
