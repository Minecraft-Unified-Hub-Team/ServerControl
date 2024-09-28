package health

import (
	"context"
)

type HealthInterface interface {
	Ping(context.Context) error
}

func NewActionService() (*HealthService, error) {
	return &HealthService{}, nil
}

type HealthService struct{}

func (hs *HealthService) Ping(ctx context.Context) error {
	return nil
}
