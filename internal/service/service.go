package service

import "github.com/AnisaForWork/user_orders/internal/config"

// Repository  holds declaration of all needed repository services used to call db level logic
type Repository interface {
}

type Service struct {
}

// NewRouter returns instane of business loggic
func NewService(repo Repository, srvCfg *config.Service) *Service {
	s := &Service{}
	return s
}
