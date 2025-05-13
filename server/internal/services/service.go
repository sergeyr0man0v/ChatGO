package services

import (
	"chatgo/server/internal/interfaces"
	"chatgo/server/internal/models"
	"time"
)

type service struct {
	models.Repository
	timeout time.Duration
}

func NewService(repository models.Repository) interfaces.Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}
