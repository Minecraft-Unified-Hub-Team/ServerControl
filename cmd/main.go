package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/action"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/config"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/health"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/server_control"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	SERVICE_CONTROL_PORT = "SERVICE_CONTROL_PORT"
)

func main() {
	value, ok := os.LookupEnv(SERVICE_CONTROL_PORT)
	if !ok {
		value = "10080"
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	logrus.WithFields(logrus.Fields{
		SERVICE_CONTROL_PORT: value,
	}).Info("tcp conn init")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", value))
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("service control init")
	actionService, err := action.NewActionService()
	if err != nil {
		logrus.Fatal(err)
	}

	healthService, err := health.NewHealthService()
	if err != nil {
		logrus.Fatal(err)
	}

	configService, err := config.NewConfigService()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("grpc init")
	grpcServer := grpc.NewServer()

	serviceControlHandler, err := server_control.NewServerControlHandler(actionService, healthService, configService)
	if err != nil {
		logrus.Fatal(err)
	}
	api.RegisterActionServer(grpcServer, serviceControlHandler)
	api.RegisterHealthServer(grpcServer, serviceControlHandler)
	api.RegisterConfigServer(grpcServer, serviceControlHandler)
	reflection.Register(grpcServer)

	logrus.Info("server init")
	grpcServer.Serve(listener)
}
