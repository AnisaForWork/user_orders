package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/AnisaForWork/user_orders/internal/config"

	_ "github.com/go-sql-driver/mysql" //nolint:blank-imports // only for sqlx
	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"
)

const (
	timeOut        = time.Second * 3
	noUser         = -1
	NumOfOrdForAvg = 50
)

var (
	ErrUniqConstrViolation = errors.New("unique constraints violation")
)

// NewMysqlD connects to Mysql with URL scheme [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func NewMysqlDB(cfg *config.DB) (*sqlx.DB, error) {
	// Mysql URL scheme [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	log.Info(fmt.Sprintf("%s:%s@tcp(%s:%s)@/%s?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Options))
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Options))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	err = db.Ping()

	return db, err
}

// Service contains db layer(postgres) as methods of this struct
type Repository struct {
	db *sqlx.DB
}

// NewRouter returns instane db layer(postgres)
func NewMysqlRepository(db *sqlx.DB) *Repository {
	repo := &Repository{
		db: db,
	}

	return repo
}

// Close for graceful shutdown
func (r *Repository) Close() error {
	return r.db.Close()
}
