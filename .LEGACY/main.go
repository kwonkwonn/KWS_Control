package main

import (
	"fmt"
	_ "os"

	api "github.com/easy-cloud-Knet/KWS_Control/api/server"
	//WorkerConn "github.com/easy-cloud-Knet/KWS_Control/api/workercont"
	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func main() {
	fmt.Println("hellot")

	//var TaskHandlersPool WorkerConn.TaskHandler
	//WorkerConn.InitWorkers(&TaskHandlersPool)
	contextStruct, err := vms.InitializeDevices("./vm/database.json")
	if err != nil {
		panic(err)
	}

	go func() {
		err := api.Server(8081, contextStruct)
		if err != nil {
			panic(err)
		}
	}()
	select {}
}
