package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// HTTPServer holds config information for http server
type HTTPServer struct {
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	Port              string
	Host              string
}

// Server contains setting for setting up Server that contains HTTP and gRPS servers
type Server struct {
	HTTP       *HTTPServer
	TimeOutSec int
}

// Configurator used to get all configurations from configured conf files
type Configurator struct {
}

// NewConfigure set paths to config.yml
func NewConfigure() (*Configurator, error) {

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read conf file: %w", err)
	}

	log.WithFields(log.Fields{
		"configFile": viper.ConfigFileUsed(),
	}).Info("config initialized")

	c := &Configurator{}

	return c, nil
}

// SrvConfig returns configuration for server
func (cfg *Configurator) HTTPSrvConfig() *HTTPServer {
	log.WithFields(log.Fields{
		"source": viper.ConfigFileUsed(),
	}).Info("reading server configuration from file")

	return &HTTPServer{
		MaxHeaderBytes:    viper.GetInt("srv.maxHeaderBytes"),
		ReadTimeout:       viper.GetDuration("srv.readTimeout"),
		WriteTimeout:      viper.GetDuration("srv.writeTimeout"),
		ReadHeaderTimeout: viper.GetDuration("srv.readHeaderTimeout"),
		Port:              viper.GetString("srv.port"),
		Host:              viper.GetString("srv.host"),
	}
}

// ServerConfig returns config for creating Server struct
func (cfg *Configurator) ServerConfig() *Server {
	httpSrv := cfg.HTTPSrvConfig()

	srv := Server{
		HTTP:       httpSrv,
		TimeOutSec: viper.GetInt("srv.timeOutSec"),
	}

	return &srv
}

type AppEnvs int

const (
	Release = iota
	Development
)

// AppENV returns application development stage
func (cfg *Configurator) AppENV() AppEnvs {
	log.WithFields(log.Fields{
		"source": ".env",
	}).Info("reading APP_ENV")

	switch os.Getenv("APP_ENV") {
	case "release":
		log.Info("application runs in release mode")

		return Release
	default:
		log.Info("application runs in Development mode")

		return Development
	}
}

// DB holds postgres config information
type DB struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	Options         string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ReconnRetry     int
	TimeWaitPerTry  time.Duration
}

// DBConfig returns configuration for postgres
func (cfg *Configurator) DBConfig() (*DB, error) {
	log.WithFields(log.Fields{
		"source1": viper.ConfigFileUsed(),
		"source2": ".env",
	}).Info("reading postgres configuration from file")

	pwd, ok := os.LookupEnv("MYSQL_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("MYSQL_PASSWORD was not found in env")
	}

	db := &DB{
		Host:            viper.GetString("mysql.host"),
		Port:            viper.GetString("mysql.port"),
		Username:        viper.GetString("mysql.username"),
		DBName:          viper.GetString("mysql.dbname"),
		Options:         viper.GetString("mysql.options"),
		Password:        pwd,
		ConnMaxLifetime: viper.GetDuration("mysql.connMaxLifetime"),
		MaxOpenConns:    viper.GetInt("mysql.maxOpenConns"),
		MaxIdleConns:    viper.GetInt("mysql.MaxIdleConns"),
		ReconnRetry:     viper.GetInt("mysql.retry"),
		TimeWaitPerTry:  viper.GetDuration("mysql.timeWaitPerTry"),
	}
	return db, nil
}

// Service holds config information all defined services
type Service struct {
}

func (cfg *Configurator) ServiceConfig() *Service {
	return &Service{}
}
