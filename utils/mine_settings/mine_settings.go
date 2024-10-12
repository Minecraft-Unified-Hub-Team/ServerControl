package mine_settings

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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
		_, err = file.Write([]byte(fmt.Sprintf("%s=%s\n", key, value)))
		if err != nil {
			return fmt.Errorf(errorFormat, err)
		}
	}

	return err
}

func ReadSettingsConfig(ctx context.Context, pathToDir string) (map[string]string, error) {
	var err error = nil
	var errorFormat string = fmt.Sprintf("mine_settings.WriteSettingsConfig(ctx, %s)", pathToDir) + ": %w"

	file, err := os.OpenFile(pathToDir+"/server.properties", os.O_RDONLY, 0660)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	defer file.Close()

	m := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && scanner.Text() != "" {
		parts := strings.Split(scanner.Text(), "=")
		if len(parts) < 2 {
			return nil, fmt.Errorf("config line bad formatted %s", scanner.Text())
		}
		key, value := parts[0], parts[1]
		m[key] = value
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return m, err
}
