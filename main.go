package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/easy-cloud-Knet/KWS_Control/structure"

	"github.com/easy-cloud-Knet/KWS_Control/api"
	"github.com/easy-cloud-Knet/KWS_Control/startup"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

func main() {
	temp1 := Find_subnet("192.168.55.")
	temp2 := Find_subnet("192.168.255.")
	temp3 := Find_subnet("192.255.255.")

	res := fmt.Sprintf("%s   /     %s      /     %s", temp1, temp2, temp3)

	fmt.Println(res)

	log := util.GetLogger()

	ctx := context.Background()

	// Redis 초기화
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

func Find_subnet(last_subnet string) string {

	value := make([]int, 3)
	var j int
	for i := 0; i < 3; i++ {
		var temp string
		for last_subnet[j] != '.' {
			temp = temp + string(last_subnet[j])
			j++
		}
		value[i], _ = strconv.Atoi(temp)
	}

	if value[2] >= 255 {
		if value[1] >= 255 {
			if value[0] >= 255 {
				return "err"
			} else {
				value[0]++
				value[1] = 0
				value[2] = 0
			}
		} else {
			value[1]++
			value[2] = 0
		}
	} else {
		value[2]++
	}

	result := fmt.Sprintf("%s.%s.%s", strconv.Itoa(value[0]), strconv.Itoa(value[1]), strconv.Itoa(value[2]))
	return result
}
