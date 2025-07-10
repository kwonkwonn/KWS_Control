package startup

import (
	"fmt"
	"os"

	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
	"gopkg.in/yaml.v3"
)

func readConfig(path string) (structure.Config, error) {
	log := util.GetLogger()

	file, err := os.Open(path)
	if err != nil {
		log.Error("failed to open config file: %v", err, true)
		return structure.Config{}, fmt.Errorf("failed to open config file: %w", err)
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
