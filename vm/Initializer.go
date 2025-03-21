package vms

import (
	//"github.com/easy-cloud-Knet/KWS_Control/config"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "gopkg.in/yaml.v3"
)

func InitializeDevices(filename string) (*ControlInfra, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &ControlInfra{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return &ControlInfra{}, fmt.Errorf("failed to read file: %v", err)
	}
	var infra ControlInfra
	infra.VMLocation = make(map[UUID]*Core)
	if err := json.Unmarshal(data, &infra); err != nil {
		return &ControlInfra{}, fmt.Errorf("failed to parse JSON: %v", err)
	}
	for i := range infra.Cores { // 인덱스로 접근하여 원본 데이터 사용
		for vmUUID := range infra.Cores[i].VMInfoIdx {
			infra.VMLocation[vmUUID] = &infra.Cores[i] // 원본 Core를 참조
		}
	}

	return &infra, nil

}
