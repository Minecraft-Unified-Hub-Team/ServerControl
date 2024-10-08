package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/avast/retry-go"
	"github.com/cucumber/godog"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker_client "github.com/docker/docker/client"
)

func NewFeatureManager(ctx context.Context) (*FeatureManager, error) {
	fm := &FeatureManager{}
	cli, err := docker_client.NewClientWithOpts(docker_client.FromEnv, docker_client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	fm.cli = cli

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", "app=server_control")

	err = RetryFunction(
		func() error {
			serverControlContainers, err := cli.ContainerList(ctx, container.ListOptions{
				Filters: filterArgs,
				All:     true,
			})
			if err != nil {
				return err
			}

			if len(serverControlContainers) > 0 {
				fm.containerId = serverControlContainers[0].ID
			} else {
				return fmt.Errorf("can't find the container")
			}

			return nil
		},
		StepOptions[FIRST_CONTAINER_WAIT].(int64),
	)

	return fm, err
}

type FeatureManager struct {
	actionServiceClient api.ActionClient
	healthServiceClient api.HealthClient

	cli         *docker_client.Client
	containerId string
	lastError   error
}

func (fm *FeatureManager) StepCleanup(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (fm *FeatureManager) serverControlIsUp(ctx context.Context) (context.Context, error) {
	err := RetryFunction(
		func() error {
			info, err := fm.cli.ContainerInspect(ctx, fm.containerId)
			if err != nil {
				return err
			}
			if info.State.Running {
				return nil
			}
			return fmt.Errorf("server isn't running")
		},
		StepOptions[DEFAULT_TIMEOUT].(int64),
	)
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	return ctx, nil
}

func (fm *FeatureManager) iInstallServer(ctx context.Context, TestServerVersion string) (context.Context, error) {
	_, err := fm.actionServiceClient.Install(context.Background(), &api.InstallRequest{Version: TestServerVersion})
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	return ctx, nil
}

func (fm *FeatureManager) iSetOption(ctx context.Context, OptionName string, OptionValue string) (context.Context, error) {
	if _, ok := StepOptions[OptionName]; !ok {
		fm.lastError = fmt.Errorf("there isn't an option with the name %s", OptionName)
		return ctx, nil
	}

	var value interface{}
	var err error
	optionType := reflect.TypeOf(StepOptions[OptionName]).String()
	switch optionType {
	case "int":
		value, err = strconv.Atoi(OptionValue)
	case "int64":
		value, err = strconv.ParseInt(OptionValue, 10, 64)
	case "float32":
		value, err = strconv.ParseFloat(OptionValue, 32)
	case "float64":
		value, err = strconv.ParseFloat(OptionValue, 64)
	case "string":
		value = OptionValue
	default:
		fm.lastError = fmt.Errorf("the variable %s has an unkown type: %s", OptionName, optionType)
		return ctx, nil
	}
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	StepOptions[OptionName] = value
	return ctx, nil
}

func (fm *FeatureManager) optionEqualTo(ctx context.Context, OptionName string, ExpectedOptionValue string) (context.Context, error) {
	if _, ok := StepOptions[OptionName]; !ok {
		fm.lastError = fmt.Errorf("there isn't an option with the name %s", OptionName)
		return ctx, nil
	}

	if fmt.Sprint(StepOptions[OptionName]) != ExpectedOptionValue {
		fm.lastError = fmt.Errorf("%s isn't equal to %s", OptionName, ExpectedOptionValue)
		return ctx, nil
	}
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

func (fm *FeatureManager) iStopServer(ctx context.Context) (context.Context, error) {
	_, err := fm.actionServiceClient.Stop(context.Background(), &api.StopRequest{})
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	return ctx, nil
}

func (fm *FeatureManager) iPingToTheServer(ctx context.Context) (context.Context, error) {
	err := RetryFunction(
		func() error {
			_, err := fm.healthServiceClient.Ping(ctx, &api.PingRequest{})
			return err
		},
		StepOptions[DEFAULT_TIMEOUT].(int64),
	)
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	return ctx, nil
}

func (fm *FeatureManager) iHaveAnError(ctx context.Context, errorDescription string) (context.Context, error) {
	defer func() {
		fm.lastError = nil
	}()

	if fm.lastError == nil {
		return ctx, fmt.Errorf("expected an error")
	}
	if fm.lastError.Error() != errorDescription {
		return ctx, fmt.Errorf("incorrect error, expected: %s, got: %s", errorDescription, fm.lastError.Error())
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
		fmt.Sprintf(":%d", StepOptions[PORT]),
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

func (fm *FeatureManager) iGetServerState(ctx context.Context, expectedStateJSON *godog.DocString) (context.Context, error) {
	jsonMap := map[string]string{}
	err := json.Unmarshal([]byte(expectedStateJSON.Content), &jsonMap)
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}

	expectedResp := &api.StateResponse{}
	expectedResp.State = api.State(api.State_value[jsonMap["State"]])

	resp, err := fm.healthServiceClient.GetState(ctx, &api.StateRequest{})
	if err != nil {
		fm.lastError = err
		return ctx, nil
	}
	if expectedResp.State != resp.State {
		fm.lastError = fmt.Errorf("get {%v} state, but {%v} state was expected", resp.State, expectedResp.State)
		return ctx, nil
	}
	return ctx, nil
}

func RetryFunction(f func() error, timeout int64) error {
	const MaxUint = ^uint(0)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	err := retry.Do(
		f,
		retry.OnRetry(func(n uint, err error) {
			logrus.Printf("%d: %s\n", n, err.Error())
		}),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Second * 3
		}),
		retry.Context(ctx),
		retry.LastErrorOnly(true),
		retry.Attempts(MaxUint),
	)
	return err
}
