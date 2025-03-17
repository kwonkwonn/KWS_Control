package vms

import (
	"fmt"
	"github.com/easy-cloud-Knet/KWS_Control/config"
	_ "gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

func InitializeDevices() (ControlInfra, error) {
	c, err := config.ReadConfig("config.yaml")
	if err != nil {
		return ControlInfra{}, err
	}

	initialContext := ControlInfra{
		Config:     c,
		Cores:      []Core{},    // 모든 코어를 관리
		AliveVM:    []*VMInfo{}, //현재 가동중인 VM의 정보
		VMLocation: map[UUID]*Core{},
	}

	for _, core := range c.Cores {
		addr := strings.Split(core, ":")
		if len(addr) != 2 {
			panic("core address should be in format ip:port")
		}

		port, err := strconv.Atoi(addr[1])
		if err != nil {
			_ = fmt.Errorf("error converting port number %w", err)
			return ControlInfra{}, err
		}
		initialContext.Cores = append(initialContext.Cores, Core{
			IP:      addr[0],
			Port:    port,
			IsAlive: true,
		})
	}

	return initialContext, nil
	// go HeartBeatSensor(InfraCon.Computers)
	// ping으로 잘 살아있는지 확인하는 놈

}
