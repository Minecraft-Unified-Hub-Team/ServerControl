package app

import (
	"context"
)

type ActionService interface {
	Start(context.Context) error
}

func NewActionService() (ActionService, error) {
	return &ActionServiceImp{}, nil
}

type ActionServiceImp struct{}

func (asi *ActionServiceImp) Start(ctx context.Context) error {
	return nil
}
