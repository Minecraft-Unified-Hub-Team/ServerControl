package server_control

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) Update(ctx context.Context, req *api.UpdateRequest) (*api.UpdateResponse, error) {
	var err error = nil

	logrus.Debug(req)

	m, err := StructToMap(req)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}
	logrus.Debug(m)

	err = sch.configService.Update(ctx, m)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.UpdateResponse{}, err
}

func StructToMap(str interface{}) (map[string]string, error) {
	typeOf := reflect.TypeOf(str)
	if typeOf.Kind() == reflect.Pointer {
		str = reflect.Indirect(reflect.ValueOf(str))
	}

	typeOf = reflect.TypeOf(str)
	if typeOf.Kind() != reflect.Struct {
		return nil, fmt.Errorf("can't reflect the fields of non-struct type %s", typeOf.Elem().Name())
	}

	m := map[string]string{}
	fields := reflect.VisibleFields(reflect.TypeOf(str))
	for _, f := range fields {
		m[f.Name] = reflect.ValueOf(str).Elem().String()
	}
	return m, nil
}
