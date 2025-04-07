package main

import (
	"fmt"
	_ "os"

	"github.com/easy-cloud-Knet/KWS_Control/api"
	"github.com/easy-cloud-Knet/KWS_Control/startup"
)

func main() {
	fmt.Println("hellot")

	contextStruct, err := startup.Initialize("./startup/vm_info.json", "config.yaml")
	if err != nil {
		panic(err)
	}

	go func() {
		err := api.Server(contextStruct.Config.Port, &contextStruct)
		if err != nil {
			panic(err)
		}
	}()
	select {}
}
