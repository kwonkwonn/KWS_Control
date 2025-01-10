package main

import (
	"fmt"
	_ "os"

	api "github.com/easy-cloud-Knet/KWS_Control/api/server"
	WorkerConn "github.com/easy-cloud-Knet/KWS_Control/api/workercont"
	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func main() {
	fmt.Println("hellot")

	var TaskHandlersPool WorkerConn.TaskHandler
	WorkerConn.InitWorkers(&TaskHandlersPool)//스레드를 초기화하는 함수
	contextStruct := vms.InitializeDevices()//VM을 정의하는 함수
	go func() {
		err := api.Server(8080, &TaskHandlersPool, &contextStruct)
		if err != nil {
			panic(err)
		}
	}()
	WorkerConn.PseudoRequestSender(&TaskHandlersPool)

	select {}
}
