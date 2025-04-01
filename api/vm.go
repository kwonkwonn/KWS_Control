package api

import (
	"net/http"
	"strconv"

	service "github.com/easy-cloud-Knet/KWS_Control/service"
	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

func Server(portNum int, contextStruct *vms.ControlInfra) error {

	http.HandleFunc("/vm", func(w http.ResponseWriter, r *http.Request) {
		if !util.CheckMethod(w, r, http.MethodPost) {
			return
		}

		err := service.CreateVM(w, r, contextStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("VM created successfully"))
	})

	err := http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
	if err != nil {
		return err
	}

	return nil
}
