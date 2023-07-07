package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AnisaForWork/user_orders/internal/config"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" //nolint:blank-imports // only for sqlx
	"github.com/jmoiron/sqlx"
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
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Options)
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	err = db.Ping()

	return db, err
}

// Service contains db layer(mysql) as methods of this struct
type Repository struct {
	db *sqlx.DB
}

// NewRouter returns instane db layer(mysql)
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

// User is db layer user model
type User struct {
	ID       int64     `json:"id" `
	Login    string    `json:"login" `
	FullName string    `json:"fullName" `
	Email    string    `json:"email" `
	Password []byte    `json:"password" `
	Created  time.Time `json:"created"`
}

// CreateUser insert use in db
func (r *Repository) CreateUser(ctx context.Context, user *User) error {

	query := "INSERT INTO users (login, email, password, fullName) values (?,?,?,?) "

	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	res, err := r.db.ExecContext(ctx, query, user.Login, user.Email, user.Password, user.FullName)

	if err != nil {
		errMsql, ok := err.(*mysql.MySQLError)
		if ok && errMsql.Number == 1062 {
			return ErrUniqConstrViolation
		}
		return err
	}

	_, err = res.LastInsertId()
	return err
}

// UserRegistered gives first user with given phone number
func (r *Repository) UserRegistered(ctx context.Context, login string, pwd []byte) error {
	query := "SELECT 1 FROM users WHERE login=? AND password=? LIMIT 1"

	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	var checker int

	return r.db.QueryRowContext(ctx, query, login, pwd).Scan(&checker)
}

type Product struct {
	Barcode  string    `json:"barcode" `
	Name     string    `json:"name" `
	Desc     string    `json:"desc" `
	Cost     int       `json:"cost" `
	UserID   int64     `json:"userId"`
	Deleted  int64     `json:"deleted"`
	Created  time.Time `json:"created"`
	FileName string    `json:"fileName"`
}
