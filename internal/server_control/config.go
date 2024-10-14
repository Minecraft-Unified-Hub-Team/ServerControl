package server_control

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"unicode"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

func (sch *ServerControlHandler) Update(ctx context.Context, req *api.UpdateConfigRequest) (*api.UpdateConfigResponse, error) {
	logrus.Debug(req)

	m := MinecraftConfigToMap(req.Config)
	logrus.Debug(m)

	err := sch.configService.WriteSettings(ctx, m)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	return &api.UpdateConfigResponse{}, err
}

func (sch *ServerControlHandler) Get(ctx context.Context, req *api.GetConfigRequest) (*api.GetConfigResponse, error) {
	logrus.Debug(req)

	m, err := sch.configService.ReadSettings(ctx)
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	config, err := MapToMinecraftConfig(m)
	if err != nil {
		return nil, err
	}

	return &api.GetConfigResponse{
		Config: config,
	}, err
}

func MinecraftConfigToMap(str interface{}) map[string]string {
	m := map[string]interface{}{}
	structs.FillMap(str, m)

	m_low := mapKeysToMinecraftStyle(m)
	removeNilFromMap(m_low)

	return mapValueToString(m_low)
}

func MapToMinecraftConfig(m map[string]string) (*api.MinecraftConfig, error) {
	config := &api.MinecraftConfig{}
	structValue := reflect.ValueOf(config).Elem()

	for key, value := range m {
		key = fromMinecraftStyle(key)
		field := structValue.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			return nil, fmt.Errorf("incorrect map %+v have not key %s. Just have %+v", m, key, field)
		}

		var interfaceValue interface{}
		switch field.Type().Elem().Kind() {
		case reflect.Int32:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to int: %v", key, err)
			}
			interfaceValue = int32(intValue)
		case reflect.String:
			interfaceValue = value
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to bool: %v", key, err)
			}
			interfaceValue = boolValue
		default:
			return nil, fmt.Errorf("unsupported field type for %s", key)
		}

		v := reflect.New(field.Type().Elem())
		v.Elem().Set(reflect.ValueOf(interfaceValue))
		field.Set(v)
	}

	return config, nil
}

func IsNilish(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return false
}

func IsPtr(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Pointer, reflect.UnsafePointer:
		return true
	}

	return false
}

func mapKeysToMinecraftStyle(m map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{}, len(m))
	for k, v := range m {
		converted[toMinecraftStyle(k)] = v
	}
	return converted
}

func toMinecraftStyle(common string) string {
	minecraftStyleBytes := []byte{}
	for i := 0; i < len(common); i++ {
		if i != 0 && common[i] > 'A' && common[i] < 'Z' {
			minecraftStyleBytes = append(minecraftStyleBytes, '-')
		}
		minecraftStyleBytes = append(minecraftStyleBytes, byte(unicode.ToLower(rune(common[i]))))
	}
	return string(minecraftStyleBytes)
}

func fromMinecraftStyle(minecraftString string) string {
	commonStringBytes := []byte{
		byte(unicode.ToUpper(rune(minecraftString[0]))),
	}
	for i := 1; i < len(minecraftString); i++ {
		if minecraftString[i] == '-' && i+1 < len(minecraftString) {
			commonStringBytes = append(commonStringBytes, byte(unicode.ToUpper(rune(minecraftString[i+1]))))
			i += 1
			continue
		}
		commonStringBytes = append(commonStringBytes, minecraftString[i])
	}
	return string(commonStringBytes)
}

func removeNilFromMap(m map[string]interface{}) {
	for k, v := range m {
		if IsNilish(v) {
			delete(m, k)
		} else if IsPtr(v) {
			m[k] = reflect.ValueOf(v).Elem()
		}
	}
}

func mapValueToString(m map[string]interface{}) map[string]string {
	str_m := map[string]string{}
	for k, v := range m {
		str_m[k] = fmt.Sprintf("%v", v)
	}
	return str_m
}
