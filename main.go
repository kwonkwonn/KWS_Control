package main

import (
	"fmt"
	_ "os"

	"github.com/easy-cloud-Knet/KWS_Control/structure"

	"github.com/easy-cloud-Knet/KWS_Control/api"
	"github.com/easy-cloud-Knet/KWS_Control/startup"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetReportCaller(true)

	log.Infof("KWS Control Server Starting...")

	contextStruct, err := startup.Initialize("./startup/vm_info.json", "config.yaml")
	if err != nil {
		log.Errorf("Failed to initialize: %v", err)
		panic(err)
	}

	printCores(contextStruct.Cores)

	go func() {
		err := api.Server(contextStruct.Config.Port, &contextStruct)
		if err != nil {
			log.Errorf("Failed to start server: %v", err)
			panic(err)
		}
	}()
	select {}
}

func printCores(cores []structure.Core) {
	for i, core := range cores {
		fmt.Printf("Core #%d: %s\n", i, core.IP)
		fmt.Printf("  * IsAlive: %t\n", core.IsAlive)
		fmt.Printf("  * FreeMemory(GiB): %.0f\n", float64(core.FreeMemory)/1024)
		fmt.Printf("  * FreeCPU: %d\n", core.FreeCPU)
		fmt.Printf("  * FreeDisk(GiB): %.0f\n", float64(core.FreeDisk)/1024)
	}
}
