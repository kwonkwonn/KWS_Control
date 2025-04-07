package startup

import (
	"fmt"
	"io"
	"os"

	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "resources/config.yaml"

func readConfig(path string) (structure.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			srcFile, err := os.Open(defaultConfigPath)
			if err != nil {
				return structure.Config{}, fmt.Errorf("failed to open default config file '%s': %w", defaultConfigPath, err)
			}
			//goland:noinspection GoUnhandledErrorResult
			defer srcFile.Close()

			destFile, err := os.Create(path)
			if err != nil {
				return structure.Config{}, fmt.Errorf("failed to create config file '%s': %w", path, err)
			}
			//goland:noinspection GoUnhandledErrorResult
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return structure.Config{}, fmt.Errorf("failed to copy default config to '%s': %w", path, err)
			}

			if err := destFile.Close(); err != nil {
				return structure.Config{}, fmt.Errorf("failed to close newly created config file '%s': %w", path, err)
			}

			file, err = os.Open(path)
			if err != nil {
				return structure.Config{}, fmt.Errorf("failed to open newly created config file '%s': %w", path, err)
			}
		} else {
			// Other error occurred when trying to open the file initially
			return structure.Config{}, fmt.Errorf("failed to open config file '%s': %w", path, err)
		}
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	var config structure.Config
	data, err := os.ReadFile(path)
	if err != nil {
		return structure.Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return structure.Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
