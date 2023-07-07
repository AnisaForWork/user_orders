package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	fmt.Println(os.Getenv("CURR_ENV"))
	switch os.Getenv("CURR_ENV") {
	case "local":

		viper.AddConfigPath("../../configs")
		viper.SetConfigName("config")
		godotenv.Load("../../.env")
	default:
		log.Info("application runs in Development mode")

		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
	}

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
	}).Info("reading mysql configuration from file")

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
	Auth    *Auth
	Product *Product
}

// Auth holds config information required for Authentication service
type Auth struct {
	PwdSec *PwdSecurity
}

// PwdSecurity holds configuration for sequring user's password
type PwdSecurity struct {
	Salt     []byte
	Times    uint32
	Memory   uint32
	Parallel uint8
	Length   uint32
}

// ServiceCfgWithConn returns configuration for service
func (cfg *Configurator) ServiceConfig() (*Service, error) {

	auth, err := cfg.AuthSrvCfg()
	if err != nil {
		return nil, err
	}

	s := &Service{
		Auth: auth,
	}
	return s, nil
}

// AuthSrvCfg returns configuration for User Authentication service
func (cfg *Configurator) AuthSrvCfg() (*Auth, error) {
	log.WithFields(log.Fields{
		"source": ".env",
	}).Info("reading auth service configurations")

	pwdSec, err := cfg.pwdSecurityCfg()
	if err != nil {
		return nil, err
	}

	s := &Auth{
		PwdSec: pwdSec,
	}
	return s, nil
}

func (cfg *Configurator) pwdSecurityCfg() (*PwdSecurity, error) {
	salt, ok := os.LookupEnv(os.Getenv("PWD_SEC_SALT"))
	if ok {
		return nil, fmt.Errorf("password security params not set")
	}

	t, err := strconv.Atoi(os.Getenv("PWD_SEC_TIMES"))
	if err != nil {
		return nil, fmt.Errorf("password security params not set")
	}

	m, err := strconv.Atoi(os.Getenv("PWD_SEC_MEMORY"))
	if err != nil {
		return nil, fmt.Errorf("password security params not set")
	}

	p, err := strconv.Atoi(os.Getenv("PWD_SEC_PARALLEL"))
	if err != nil {
		return nil, fmt.Errorf("password security params not set")
	}

	l, err := strconv.Atoi(os.Getenv("PWD_SEC_LENGTH"))
	if err != nil {
		return nil, fmt.Errorf("password security params not set")
	}

	pwd := &PwdSecurity{
		Salt:     []byte(salt),
		Times:    uint32(t),
		Memory:   uint32(m),
		Parallel: uint8(p),
		Length:   uint32(l),
	}
	return pwd, nil
}

// Auth holds config information required for Authentication service
type Product struct {
	PathToCheckDir string
	PathToTemplate string
	TemplateName   string
	FontName       string
	FontFileName   string
	TimeFormat     string
	TemplateW      float64
	TemplateH      float64
}

// ProductConfig returns configuration for product service
func (cfg *Configurator) ProductConfig() *Product {
	log.WithFields(log.Fields{
		"source1": viper.ConfigFileUsed(),
	}).Info("reading product service configuration from file")

	pr := &Product{
		PathToCheckDir: viper.GetString("product.pathToCheckDir"),
		PathToTemplate: viper.GetString("product.pathToTemplate"),
		TemplateName:   viper.GetString("product.templateName"),
		FontName:       viper.GetString("product.fontName"),
		FontFileName:   viper.GetString("product.fontFileName"),
		TimeFormat:     viper.GetString("product.timeFormat"),
		TemplateW:      viper.GetFloat64("product.templateW"),
		TemplateH:      viper.GetFloat64("product.templateH"),
	}
	return pr
}

type JWTProvider struct {
	Host         string
	Port         int
	Timeout      time.Duration
	Retry        int
	TimeoutRetry time.Duration
}

// JWTProviderConfig returns configuration for jwt generator
func (cfg *Configurator) JWTProviderConfig() *JWTProvider {
	log.WithFields(log.Fields{
		"source1": viper.ConfigFileUsed(),
	}).Info("reading token generator configuration from file")

	jp := &JWTProvider{
		Host:         viper.GetString("tokenGen.host"),
		Port:         viper.GetInt("tokenGen.port"),
		Timeout:      viper.GetDuration("tokenGen.timeout"),
		Retry:        viper.GetInt("tokenGen.retry"),
		TimeoutRetry: viper.GetDuration("tokenGen.timeoutRetry"),
	}
	return jp
}
