package main

import (
	_ "os"

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

	go func() {
		err := api.Server(contextStruct.Config.Port, &contextStruct)
		if err != nil {
			log.Errorf("Failed to start server: %v", err)
			panic(err)
		}
	}()
	select {}
}
