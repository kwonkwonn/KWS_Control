package vms

import (
	"encoding/json"
	"fmt"
	// "github.com/easy-cloud-Knet/KWS_Control/config"
	"io"
	"os"
	_ "gopkg.in/yaml.v3"
	"strconv"
	"strings"
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
		
		addr := strings.Split(infra.Cores[i].IP+":"+strconv.Itoa(infra.Cores[i].Port), ":")
		if len(addr) != 2 {
			panic("core address should be in format ip:port")
		}

		port, err := strconv.Atoi(addr[1])
		if err != nil {
			_ = fmt.Errorf("error converting port number %w", err)
			return &ControlInfra{}, err
		}
		
		infra.Cores[i].IP = addr[0]
		infra.Cores[i].Port = port
		infra.Cores[i].IsAlive = true
	}

	return &infra, nil
	// go HeartBeatSensor(InfraCon.Computers)
	// ping으로 잘 살아있는지 확인하는 놈
}