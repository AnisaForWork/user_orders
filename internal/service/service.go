package service

import (
	"github.com/AnisaForWork/user_orders/internal/config"
	"github.com/AnisaForWork/user_orders/internal/provider/token"
	"github.com/AnisaForWork/user_orders/internal/service/auth"
)

// Repository  holds declaration of all needed repository services used to call db level logic
type Repository interface {
	auth.Repository
}

type Service struct {
	*auth.AService
}

// NewRouter returns instane of business loggic
func NewService(repo Repository, provider *token.JWTProvider, srvCfg *config.Service) *Service {
	a := auth.NewService(repo, provider, srvCfg.Auth)
	s := &Service{
		AService: a,
	}
	return s
}
