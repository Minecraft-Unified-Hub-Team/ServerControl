package mine_settings

import (
	"context"
	"fmt"
	"os"
)

func WriteSettingsConfig(ctx context.Context, pathToDir string, settings map[string]string) error {
	var err error = nil
	var errorFormat string = fmt.Sprintf("mine_settings.WriteSettingsConfig(ctx, %s)", pathToDir) + ": %w"

	file, err := os.OpenFile(pathToDir+"/server.properties", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	defer file.Close()

	for key, value := range settings {
		_, err = file.Write([]byte(fmt.Sprintf("%s:%s\n", key, value)))
		if err != nil {
			return fmt.Errorf(errorFormat, err)
		}
	}

	return err
}
