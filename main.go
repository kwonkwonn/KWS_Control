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
	WorkerConn.InitWorkers(&TaskHandlersPool)
	contextStruct, err := vms.InitializeDevices()
	if err != nil {
		panic(err)
	}

	go func() {
		err := api.Server(8080, &TaskHandlersPool, &contextStruct)
		//wg.Wait()
		//api.Unlock()
		if err != nil {
			panic(err)
		}
	}()

	//api.Done()
	//WorkerConn.PseudoRequestSender(&TaskHandlersPool)

	select {}
}
