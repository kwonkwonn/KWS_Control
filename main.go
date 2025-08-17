package main

import (
	"context"
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Control/structure"

	"github.com/easy-cloud-Knet/KWS_Control/api"
	"github.com/easy-cloud-Knet/KWS_Control/startup"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

func main() {
	log := util.GetLogger()

	ctx := context.Background()

	//Redis 초기화
	rdb, err := startup.InitializeRedis(ctx)
	if err != nil {
		log.Error("Failed to initialize Redis: %v", err, true)
		panic(err)
	}

	log.Info("KWS Control Server Starting...", true)

	contextStruct, err := startup.InitializeCoreData("config.yaml")
	if err != nil {
		log.Error("Failed to initialize: %v", err, true)
		panic(err)
	}
	printCores(contextStruct.Cores)
	// cmsClient := service.NewCmsClient()
	// addrResp := cmsClient.NewCmsSubnet("20.20.22.")
	// fmt.Printf("%s\n", addrResp.IP)
	// fmt.Printf("%s\n", addrResp.MacAddr)
	// fmt.Printf("%s\n", addrResp.SdnUUID)

	go func() {
		err := api.Server(contextStruct.Config.Port, &contextStruct, rdb)
		if err != nil {
			log.Error("Failed to start server: %v", err, true)
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
