package mysql

import (
	"context"
	"database/sql"
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
	ErrNoRows              = errors.New("no rows in table returned")
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
	ID        int64     `json:"id" `
	Login     string    `json:"login" `
	FullName  string    `json:"fullName" `
	Email     string    `json:"email" `
	Password  []byte    `json:"password" `
	Created   time.Time `json:"created"`
	FileNames []string  `json:"fileNames"`
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
	Barcode string    `json:"barcode" `
	Name    string    `json:"name" `
	Descr   string    `json:"descr" `
	Cost    int       `json:"cost" `
	UserID  int64     `json:"userId"`
	Deleted bool      `json:"deleted"`
	Created time.Time `json:"created"`
}

func (r *Repository) Create(ctx context.Context, pr Product, login string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut*3)
	defer cancel()

	query := "INSERT INTO products (barcode, name, descr, cost,userId) values (?,?,?,?,?) "
	selector := "SELECT id FROM users WHERE login = ?"

	var tx *sqlx.Tx
	tx, err = r.db.BeginTxx(ctx, &sql.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return err
	}
	var id int64

	err = tx.QueryRow(selector, login).Scan(&id)

	if err != nil {
		return ErrNoRows
	}

	res, err := r.db.ExecContext(ctx, query, pr.Barcode, pr.Name, pr.Descr, pr.Cost, id)

	if err != nil {
		errMsql, ok := err.(*mysql.MySQLError)
		if ok && errMsql.Number == 1062 {
			return ErrUniqConstrViolation
		}
		return err
	}

	_, err = res.LastInsertId()

	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (r *Repository) UserProducts(ctx context.Context, amount int, offset int, usrID string) ([]Product, error) {

	query := `SELECT barcode, name, cost FROM products 
				JOIN users ON users.id=products.userId AND users.login=? 
				ORDER BY products.created 
				LIMIT ? OFFSET ?`

	ctx, cancel := context.WithTimeout(ctx, timeOut*3)
	defer cancel()

	prods := []Product{}

	err := r.db.SelectContext(ctx, &prods, query, usrID, amount, offset)

	return prods, err
}

func (r *Repository) UserProduct(ctx context.Context, barcode string, login string) (pr *Product, prchecks []string, err error) {
	queryPr := `SELECT barcode, name, cost, descr,products.created FROM products 
					JOIN users ON users.id=products.userId AND users.login=? 
					WHERE barcode=? AND deleted=FALSE 
					LIMIT 1`
	queryCh := `SELECT filename FROM prchecks 
					WHERE prchecks.barcode=? 
					LIMIT 100`

	ctx, cancel := context.WithTimeout(ctx, timeOut*3)
	defer cancel()

	var tx *sqlx.Tx
	tx, err = r.db.BeginTxx(ctx, &sql.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	pr = &Product{}
	err = r.db.GetContext(ctx, pr, queryPr, login, barcode)
	if err != nil {
		return nil, nil, ErrNoRows
	}

	prchecks = []string{}
	err = r.db.SelectContext(ctx, &prchecks, queryCh, barcode)
	if err != nil {
		return nil, nil, ErrNoRows
	}

	err = tx.Commit()

	if err != nil {
		return nil, nil, err
	}
	return pr, prchecks, err
}
func (r *Repository) Delete(ctx context.Context, barcode string, login string) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	query := `UPDATE products SET deleted=TRUE 
				JOIN users ON users.id=products.userId AND users.login=? 
				WHERE barcode=? `

	res, err := r.db.ExecContext(ctx, query, login, barcode)

	if err != nil {
		return err
	}

	rowC, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowC == 0 {
		return ErrNoRows
	}

	return nil
}
func (r *Repository) ProductInfoForCheck(ctx context.Context, barcode string, login string) (*Product, error) {
	queryPr := `SELECT barcode, name, cost FROM products 
					JOIN users ON users.id=products.userId AND users.login=? 
					WHERE barcode=? AND deleted=FALSE 
					LIMIT 1`
	ctx, cancel := context.WithTimeout(ctx, timeOut*3)
	defer cancel()

	var pr Product

	err := r.db.GetContext(ctx, &pr, queryPr, login, barcode)
	if err != nil {
		return nil, ErrNoRows
	}

	return &pr, nil
}
func (r *Repository) UpdateCheckInfo(ctx context.Context, filename string, barcode string, login string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut*3)
	defer cancel()

	query := "INSERT INTO prchecks (filename,barcode) values (?,?) "
	selector := `SELECT 1 FROM products 
					JOIN users ON users.id=products.userId AND users.login=? 
					WHERE barcode=? AND deleted=FALSE 
					LIMIT 1`

	var tx *sqlx.Tx
	tx, err = r.db.BeginTxx(ctx, &sql.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return err
	}
	var checkUser int

	err = tx.QueryRow(selector, login, barcode).Scan(&checkUser)

	if err != nil {
		return ErrNoRows
	}

	res, err := r.db.ExecContext(ctx, query, filename, barcode)

	if err != nil {
		errMsql, ok := err.(*mysql.MySQLError)
		if ok && errMsql.Number == 1062 {
			return ErrUniqConstrViolation
		}
		return err
	}

	_, err = res.LastInsertId()

	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
func (r *Repository) CheckOwnership(ctx context.Context, filename string, login string) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	query := `SELECT 1 FROM prchecks 
				JOIN products ON prchecks.barcode=products.barcode 
				JOIN users ON users.id=products.userId AND users.login=? 
				WHERE prchecks.filename=? AND products.deleted=FALSE LIMIT 1`

	var res int
	err := r.db.QueryRowContext(ctx, query, login, filename).Scan(&res)

	if err != nil {
		return err
	}

	if res != 1 {
		return ErrNoRows
	}

	return nil
}
