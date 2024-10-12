package server_control

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Minecraft-Unified-Hub-Team/ServerControl/internal/api"
)

func TestMinecraftConfigToMap(t *testing.T) {
	hardcode := true
	max_players := int32(555)
	motd := "It's just another Minecraft server."

	req := &api.UpdateConfigRequest{
		Config: &api.MinecraftConfig{
			Hardcore:   &hardcode,
			MaxPlayers: &max_players,
			Motd:       &motd,
		},
	}

	exp := map[string]string{
		"max-players": "555",
		"hardcore":    "true",
		"motd":        "It's just another Minecraft server.",
	}
	res := MinecraftConfigToMap(req.Config)
	if err := mapEqual(exp, res); err != nil {
		t.Fatal(exp, res, err)
	}
}

func TestMapToMinecraftConfig(t *testing.T) {
	hardcode := true
	max_players := int32(555)
	motd := "It's just another Minecraft server."

	exp := &api.UpdateConfigRequest{
		Config: &api.MinecraftConfig{
			Hardcore:   &hardcode,
			MaxPlayers: &max_players,
			Motd:       &motd,
		},
	}

	m := map[string]string{
		"max-players": "555",
		"hardcore":    "true",
		"motd":        "It's just another Minecraft server.",
	}
	res, err := MapToMinecraftConfig(m)
	if err != nil {
		t.Fatal(exp.Config, res, err)
	}

	if ok := reflect.DeepEqual(exp.Config, res); !ok {
		t.Fatal(exp, res, ok)
	}
}

func mapEqual(exp, res map[string]string) error {
	for k, v := range exp {
		if res[k] != v {
			return fmt.Errorf("key: %s, exp: %s, res: %s", k, exp[k], res[k])
		}
	}

	for k, v := range res {
		if exp[k] != v {
			return fmt.Errorf("key: %s, exp: %s, res: %s", k, exp[k], res[k])
		}
	}

	return nil
}
