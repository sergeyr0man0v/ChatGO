package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"time"
)

type service struct {
	models.Repository
	timeout time.Duration
	Config
}

type Config struct {
	secretKey  string `yaml:"JWTKey"`
	encryptKey []byte `yaml:"encryptKey"`
}

func NewService(repository models.Repository, config *Config) interfaces.Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
		*config,
	}
}
