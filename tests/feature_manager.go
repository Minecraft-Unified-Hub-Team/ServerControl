package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"
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
	actionServiceClient api.ActionClient
	healthServiceClient api.HealthClient

	lastError error
}

func (fm *FeatureManager) StepCleanup(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (fm *FeatureManager) iStartServer(ctx context.Context) (context.Context, error) {
	_, err := fm.actionServiceClient.Start(context.Background(), &api.StartRequest{})
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	return ctx, nil
}

func (fm *FeatureManager) iPingToTheServer(ctx context.Context) (context.Context, error) {
	err := retry.Do(
		func() error {
			_, err := fm.healthServiceClient.Ping(ctx, &api.PingRequest{})
			return err
		},
		retry.OnRetry(func(n uint, err error) {
			logrus.Printf("%d: %s\n", n, err.Error())
			fmt.Printf("%d: %s\n", n, err.Error())
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Second * time.Duration(n)
		}),
		retry.Attempts(5),
	)
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	return ctx, nil
}

func (fm *FeatureManager) iHaveAnError(ctx context.Context) (context.Context, error) {
	if fm.lastError == nil {
		return ctx, fmt.Errorf("expected an error")
	}
	return ctx, nil
}

func (fm *FeatureManager) iHaveNoErrors(ctx context.Context) (context.Context, error) {
	if fm.lastError != nil {
		return ctx, fm.lastError
	}
	return ctx, nil
}

func (fm *FeatureManager) iConnectToServiceControl(ctx context.Context) (context.Context, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", 10080),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}

	fm.actionServiceClient = api.NewActionClient(conn)
	fm.healthServiceClient = api.NewHealthClient(conn)

	return ctx, nil
}