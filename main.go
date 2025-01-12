package main

import (
	"context"
	"fmt"
	_ "os"

	api "github.com/easy-cloud-Knet/KWS_Control/api/server"
	WorkerConn "github.com/easy-cloud-Knet/KWS_Control/api/workercont"
	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func main() {
	fmt.Println("hellot")

	ctx := context.Background()

	var TaskHandlersPool WorkerConn.TaskHandler
	WorkerConn.InitWorkers(&TaskHandlersPool)
	contextStruct := vms.InitializeDevices(ctx)
	go func() {
		err := api.Server(8080, &TaskHandlersPool, &contextStruct)
		if err != nil {
			panic(err)
		}
	}()
	WorkerConn.PseudoRequestSender(&TaskHandlersPool)

	select {}
}
