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
	"github.com/AnisaForWork/user_orders/internal/repository/mysql"
	"github.com/AnisaForWork/user_orders/internal/server"
	"github.com/AnisaForWork/user_orders/internal/service"
	"github.com/AnisaForWork/user_orders/migration"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.NewConfigure()

	if err != nil {
		log.WithFields(log.Fields{
			"place": "system(main)",
		}).WithError(err).Panic("Config initialization failed")
	}

	lgf := log.Fields{}
	isRelease := cfg.AppENV() != config.Development
	if !isRelease {
		lgf["APP_ENV"] = "development"
	} else {
		lgf["APP_ENV"] = "production"
	}

	log.WithFields(lgf).Info("Read application environment")

	migrator := migration.NewMigratory()

	repo, err := initDBwithRetry(cfg, migrator)
	if err != nil {
		log.WithFields(log.Fields{
			"place": "system(main)",
		}).WithError(err).Panic("Initialization of postgres repo failed")
	}
	defer repo.Close()

	servCfg := cfg.ServiceConfig()
	serv := service.NewService(repo, servCfg)

	hndl := handler.NewRouter(serv, isRelease)

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

func initDBwithRetry(cfg *config.Configurator, migr *migration.Migratory) (repo *mysql.Repository, err error) {

	for i := 0; i < 10; i++ {
		repo, err = initDB(cfg, migr)
		if err == nil {
			break
		}
		log.WithFields(log.Fields{
			"retryNumber": i,
		}).WithError(err).Warn("Could not initialize DB, trying to retry")
		time.Sleep(time.Second * 3)
	}

	return repo, err
}

// initDB reads db configuration, connects to postgres, does migration if needed
// returns struct that implements logic.Repository interface or error
func initDB(cfg *config.Configurator, migr *migration.Migratory) (repo *mysql.Repository, err error) {
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
	db, err := mysql.NewMysqlDB(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("mysql failed to initialize with error %w", err)
	}

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
