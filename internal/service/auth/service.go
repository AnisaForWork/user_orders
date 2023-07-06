package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AnisaForWork/user_orders/internal/config"
	"github.com/AnisaForWork/user_orders/internal/repository/mysql"

	"golang.org/x/crypto/argon2"
)

// Repository used to call db level logic
type Repository interface {
	CreateUser(context.Context, *mysql.User) (err error)
	UserRegistered(context.Context, string, []byte) error
}

type Provider interface {
	UserJWTToken(login string) (string, error)
}

// PwdSecurity holds cryptography configurations to sequre user's password
type PwdSecurity struct {
	Salt     []byte
	Times    uint32
	Memory   uint32
	Parallel uint8
	Length   uint32
}

// AService struct implements auth service functionality
type AService struct {
	Repo     Repository
	TokenGen Provider
	PwdSec   PwdSecurity
}

func NewService(repo Repository, prvd Provider, srvCfg *config.Auth) *AService {
	pwdSec := PwdSecurity{
		Salt:     srvCfg.PwdSec.Salt,
		Times:    srvCfg.PwdSec.Times,
		Memory:   srvCfg.PwdSec.Memory,
		Parallel: srvCfg.PwdSec.Parallel,
		Length:   srvCfg.PwdSec.Length,
	}

	s := &AService{
		Repo:     repo,
		TokenGen: prvd,
		PwdSec:   pwdSec,
	}
	return s
}

var (
	ErrNoSuchUser      = errors.New("user with such credential doesn't exist")
	ErrInvalidToken    = errors.New("systme did not create this token")
	ErrExpiredToken    = errors.New("token expired")
	ErrUserAlredyExist = errors.New("user with such credentials already registered")
	ErrUsrDelted       = errors.New("user deleted account")
)

// User is service level user model
type User struct {
	ID       int64
	Login    string
	FullName string
	Email    string
	Password string
}

// SingUp register new user(with uniq email and phone) by adding them to db
func (s *AService) SingUp(ctx context.Context, usr *User) error {
	pwd := argon2.IDKey([]byte(usr.Password), s.PwdSec.Salt, s.PwdSec.Times, s.PwdSec.Memory, s.PwdSec.Parallel, s.PwdSec.Length)

	user := mysql.User{
		Login:    usr.Login,
		Email:    usr.Email,
		FullName: usr.FullName,
		Password: pwd,
	}

	err := s.Repo.CreateUser(ctx, &user)
	if errors.Is(err, mysql.ErrUniqConstrViolation) {
		return ErrUserAlredyExist
	}

	return err
}

// Auth is service level authentication info model
type Auth struct {
	Login    string
	Password string
}

// SingIn authenticate using given credentials
// returns acess+resfresh tokens if credentials valid
func (s *AService) SingIn(ctx context.Context, auth Auth) (string, error) {
	pwd := argon2.IDKey([]byte(auth.Password), s.PwdSec.Salt, s.PwdSec.Times, s.PwdSec.Memory, s.PwdSec.Parallel, s.PwdSec.Length)

	err := s.Repo.UserRegistered(ctx, auth.Login, pwd)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoSuchUser
	}

	return s.TokenGen.UserJWTToken(auth.Login)
}
