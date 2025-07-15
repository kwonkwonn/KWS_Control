package startup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
	"gopkg.in/yaml.v3"
)

func readConfig(path string) (structure.Config, error) {
	log := util.GetLogger()

	// 받아온 path 읽기
	config, err := tryReadConfig(path)
	if err == nil {
		log.DebugInfo("Successfully read config from: %s", path)
		return config, nil
	}

	log.Warn("Failed to read config from %s: %v", path, err)

	fallbackPath := filepath.Join("resources", "config.yaml")
	log.Info("Attempting to read config from fallback path: %s", fallbackPath)

	// 받아온 path에 없으면 resources/에 있는 config.yaml 가져오기
	config, fallbackErr := tryReadConfig(fallbackPath)
	if fallbackErr != nil {
		return structure.Config{}, fmt.Errorf("both original path (%s) and fallback path (%s) failed. Original error: %w, Fallback error: %v", path, fallbackPath, err, fallbackErr)
	}

	log.Info("Successfully read config from fallback path: %s", fallbackPath)
	return config, nil
}

// tryReadConfig attempts to read and parse a config file from the given path
func tryReadConfig(path string) (structure.Config, error) {
	file, err := os.Open(path)
	if err != nil {
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
