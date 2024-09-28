package action

import (
	"context"
)

type ActionInterface interface {
	Start(context.Context) error
}

func NewActionService() (*ActionService, error) {
	return &ActionService{}, nil
}

type ActionService struct{}

func (as *ActionService) Start(ctx context.Context) error {
	return nil
}
