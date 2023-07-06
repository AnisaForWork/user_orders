package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnisaForWork/user_orders/internal/config"
	"github.com/AnisaForWork/user_orders/internal/handler"
	"github.com/AnisaForWork/user_orders/internal/provider/token"
	"github.com/AnisaForWork/user_orders/internal/repository/mysql"
	"github.com/AnisaForWork/user_orders/internal/server"
	"github.com/AnisaForWork/user_orders/internal/service"
	"github.com/AnisaForWork/user_orders/migration"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type ServicesAndProviders struct {
	*service.Service
	*token.JWTProvider
}

func main() {
	cfg, err := config.NewConfigure()

	if err != nil {
		log.WithFields(log.Fields{
			"place": "system(main)",
		}).WithError(err).Panic("Config initialization failed")
	}

	log.Info("Configuration provider initialized")

	lgf := log.Fields{}
	isRelease := cfg.AppENV() != config.Development
	if !isRelease {
		lgf["APP_ENV"] = "development"
	} else {
		lgf["APP_ENV"] = "production"
	}

	log.WithFields(lgf).Info("Got application environment")

	initLogger(cfg, isRelease)

	migrator := migration.NewMigratory()

	repo, err := initDBwithRetry(cfg, migrator)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "system(main)",
		}).WithError(err).Panic("Initialization of mysql repo failed")
	}
	defer repo.Close()

	providerCfg := cfg.JWTProviderConfig()
	provider, err := token.NewJWTProvider(*providerCfg)
	if err != nil {
		log.WithFields(log.Fields{"place": "system(main)"}).WithError(err).Panic("Failed to connect to token generator")
	}

	servCfg, err := cfg.ServiceConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"place": "system(main)",
		}).WithError(err).Panic("Cant access service config")
	}
	serv := service.NewService(repo, provider, servCfg)

	srvWPrv := ServicesAndProviders{
		serv,
		provider,
	}

	hndl := handler.NewRouter(srvWPrv, isRelease)

	srvCfg := cfg.ServerConfig()
	srv := server.NewServer(srvCfg, hndl)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := srv.Run(ctx); err != nil {
			log.WithFields(log.Fields{"place": "system(main)"}).WithError(err).Error("Server failed during run")
		}
	}() // adding more code here can cause sync problems(programm can stop befor end of goroutine)

	shutDown := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := <-shutDown
	log.WithFields(log.Fields{
		"signal": s,
	}).Info("System received signal for shutdown,shutting down server")

	cancel()

}

func initLogger(cfg *config.Configurator, isRelease bool) {
	log.SetOutput(os.Stdout)

	if !isRelease {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetLevel(log.WarnLevel)
		log.SetFormatter(&log.JSONFormatter{})
	}
}

// initDB reads db configuration, connects to mysql, does migration if needed
// returns struct that implements logic.Repository interface or error
func initDBwithRetry(cfg *config.Configurator, migr *migration.Migratory) (repo *mysql.Repository, err error) {
	log.WithFields(log.Fields{
		"cfg for": "mysql",
	}).Info("Getting configs for db")

	dbCfg, err := cfg.DBConfig()
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"dbType": "mysql",
	}).Info("Getting db")

	log.WithFields(log.Fields{"conf": dbCfg}).Info("DB config")

	var db *sqlx.DB
	for i := 0; i < dbCfg.ReconnRetry; i++ {

		db, err = mysql.NewMysqlDB(dbCfg)
		if err != nil {
			return nil, fmt.Errorf("mysql failed to initialize with error %w", err)
		}

		if err == nil {

			log.WithFields(log.Fields{
				"dbType": "mysql",
			}).Info("Migrating db")

			if err = migr.Migrate(db); err != nil {
				return nil, fmt.Errorf("mysql failed to migrate with error: %w", err)
			}

			repo = mysql.NewMysqlRepository(db)

			log.WithFields(log.Fields{}).Info("Repository loaded")

			return repo, nil
		}

		log.WithFields(log.Fields{
			"retryNumber": i,
		}).WithError(err).Warn("Could not initialize DB, trying to retry")

		time.Sleep(dbCfg.TimeWaitPerTry)
	}

	return repo, err
}
