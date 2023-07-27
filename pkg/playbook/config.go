package playbook

// ConfigWorkspace is ...
type ConfigWorkspace struct {
	ConfigBasedate       ConfigBasedate    `toml:"database"`
	ConfigMail           ConfigMail        `toml:"mail"`
	URLConfig            URLConfig         `toml:"url"`
	MongoConfig          MongoConfig       `toml:"mongo"`
	PluginConfig         PluginConfig      `toml:"plugin"`
	RedisConfig          RedisConfig       `toml:"redis"`
	RedisSessionConfig   RedisConfig       `toml:"redis_session"`
	Env                  map[string]string `toml:"env"`
	HttpsEngineConfig    HttpsConfig       `toml:"https_engine"`
	HttpsDesingnerConfig HttpsConfig       `toml:"https_designer"`
}
type HttpsConfig struct {
	Enable      bool    `tom:"enable"`
	Cert        string  `tom:"cert"`
	Key         string  `tom:"key"`
	Address     string  `tom:"address"`
	Description string  `tom:"description"`
	User        *string `tom:"user"`
	Password    *string `tom:"password"`
}

type RedisSessionConfig struct {
	Host     string `tom:"host"`
	Password string `tom:"password"`
}

type RedisConfig struct {
	Host              string `tom:"host"`
	Password          string `tom:"password"`
	MaxConnectionPool int    `tom:"max_connection_pool"`
}

type MongoConfig struct {
	URL string `tom:"url"`
}

type PluginConfig struct {
	Plugins []string `toml:"plugins"`
}

type URLConfig struct {
	URLBase string `toml:"url_base"`
}

// ConfigBasedate is ...
type ConfigBasedate struct {
	DatabaseURL    string `toml:"url"`
	DatabaseDriver string `toml:"driver"`
	DatabaseInit   string `toml:"init"`
}

// ConfigMail is ...
type ConfigMail struct {
	MailSMTP     string `toml:"smtp"`
	MailSMTPPort string `toml:"port"`
	MailFrom     string `toml:"from"`
	MailPassword string `toml:"password"`
}
