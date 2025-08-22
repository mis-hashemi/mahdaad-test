package order_service

import "github.com/mis-hashemi/mahdaad-test/pkg/logger"

type Service struct {
	log *logger.Logger
}

func New(log *logger.Logger) *Service {
	return &Service{log: log}
}
