package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/app"
	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/impl"
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
	actionService, err := app.NewActionService()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("grpc init")
	grpcServer := grpc.NewServer()

	serviceControlHandler, err := impl.NewServiceControlHandler(actionService)
	if err != nil {
		logrus.Fatal(err)
	}
	api.RegisterActionServiceServer(grpcServer, serviceControlHandler)
	reflection.Register(grpcServer)

	logrus.Info("server init")
	grpcServer.Serve(listener)
}
