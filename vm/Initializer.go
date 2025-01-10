package vms

import (
	"fmt"

	_ "gopkg.in/yaml.v3"
)

func InitializeDevices() InfraContext {
	initialContext := InfraContext{
		Computers:         []Computer{}, //컴퓨터 목록 관리
		VMPoolUnallocated: []*VM{}, //아직 할당되지 않은 VM을 가리키는 포인터
		VMPoolAllocated:   []*VM{}, //할당된 VM을 가리키는 포인터
		VMPool:            map[UUID]*VM{},
	}
	// config 파일이나 데이터베이스에서 읽어와야 함.
	// 현재 연결되어 있는 컴퓨터들의 간단한 정보,

	COM1 := Computer{
		Name:    "worker1",
		IP:      "192.168.64.5",
		IsAlive: false,
	}
	COM2 := Computer{
		Name:    "worker2",
		IP:      "192.168.64.6",
		IsAlive: false,
	}

	//COM1과 COM2를 initialContext.Computers에 정의
	initialContext.Computers = append(initialContext.Computers, COM1, COM2)
	fmt.Println("hellot1")
	initialContext.UpdateList()//???????
	fmt.Println("hellot2")

	return initialContext
	// go HeartBeatSensor(InfraCon.Computers)
	// ping으로 잘 살아있는지 확인하는 놈

}
