package tests

import (
	"context"
	"fmt"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	PORT = 10080
)

func NewFeatureManager() (*FeatureManager, error) {
	return &FeatureManager{}, nil
}

type FeatureManager struct {
	actionServiceClient api.ActionServiceClient
}

func (fm *FeatureManager) StepCleanup(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (fm *FeatureManager) StepIStartServer(ctx context.Context) (context.Context, error) {
	_, err := fm.actionServiceClient.Start(context.Background(), &api.StartRequest{})
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (fm *FeatureManager) StepIConnectToServiceControl(ctx context.Context) (context.Context, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", 10080),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	fm.actionServiceClient = api.NewActionServiceClient(conn)
	return ctx, nil
}
