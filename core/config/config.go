package config

import (
	"fmt"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	doOnce    sync.Once
	AppConfig Configuration
)

type Server struct {
	TZ   string `envconfig:"TZ" default:"UTC"`
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"8080"`
}

func (s Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Mysql struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
	DBName   string `envconfig:"DBNAME"`
}

func (m Mysql) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DBName,
	)
}

type Database struct {
	Mysql    Mysql `envconfig:"MYSQL"`
	LogLevel int   `envconfig:"LOG_LEVEL" default:"4"`
}

type App struct {
	Environment  string   `envconfig:"ENVIRONMENT"`
	EnableCors   bool     `envconfig:"ENABLE_CORS" default:"false"`
	AllowOrigins []string `envconfig:"ALLOWORIGINS" default:"*"`
}

type MigrateConfig struct {
	AutoMigrate bool `envconfig:"AUTO" default:"false"`
}

type JWTConfig struct {
	Key                string `envconfig:"KEY"`
	AccessTokenExpiry  int    `envconfig:"ACCESS_TOKEN_EXPIRY" default:"86400"`   // default 1d
	RefreshTokenExpiry int    `envconfig:"REFRESH_TOKEN_EXPIRY" default:"604800"` // default 7d
}

type Configuration struct {
	LoggerLevel string   `envconfig:"LOG_LEVEL" default:"error"`
	Server      Server   `envconfig:"SERVER"`
	Database    Database `envconfig:"DATABASE"`
	App         App      `envconfig:"APP"`

	MigrateConfig MigrateConfig `envconfig:"MIGRATE"`

	JWT JWTConfig `envconfig:"JWT"`
}

const NAMESPACE = "LOGISTICS_APP"

func Get() Configuration {
	var configuration Configuration

	envconfig.Process(NAMESPACE, &configuration)

	return configuration
}

func init() {
	doOnce.Do(func() {
		AppConfig = Get()
	})
}
