package api

import (
	"fmt"
	"net/http"
	"strconv"

	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
)

type handlerContext struct {
	context *vms.ControlContext
}

func Server(portNum int, contextStruct *vms.ControlContext) error {
	h := handlerContext{
		context: contextStruct,
	}
	http.HandleFunc("POST /vm", h.createVm)

	fmt.Printf("Running server on port %d\n", portNum)
	err := http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
	if err != nil {
		return err
	}

	return nil
}
