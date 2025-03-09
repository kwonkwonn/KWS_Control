package vms

import (
	"github.com/easy-cloud-Knet/KWS_Control/config"
	_ "gopkg.in/yaml.v3"
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
	// config 파일이나 데이터베이스에서 읽어와야 함.
	// 현재 연결되어 있는 컴퓨터들의 간단한 정보,

	Core1 := Core{
		IP:      "192.168.64.5",
		IsAlive: false,
	}
	Core2 := Core{
		IP:      "192.168.64.6",
		IsAlive: false,
	}

	//COM1과 COM2를 initialContext.Computers에 정의
	initialContext.Cores = append(initialContext.Cores, Core1, Core2)

	return initialContext, nil
	// go HeartBeatSensor(InfraCon.Computers)
	// ping으로 잘 살아있는지 확인하는 놈

}
