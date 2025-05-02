package server

import (
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"

	WorkerCont "github.com/easy-cloud-Knet/KWS_Control/api/workercont"
	"github.com/easy-cloud-Knet/KWS_Control/util"
	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func Server(portNum int, contextStruct *vms.ControlInfra) error {
	// main server와 통신하기 위한 http 서버
	// gin.DefaultWriter = io.Discard
	// http.HandleFunc("Get /getStatus", func(w http.ResponsseWriter, r *http.Request) {

	// 	if r.Method != http.MethodGet {
	// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 		return
	// 	}
	// 	workerControl := &WorkerCont.TaskControlGetStatus{
	// 		ResultChan: make(chan WorkerCont.TaskExecutionResult),
	// 	}
	// 	resultChannel := workerControl.ResultChan
	// 	defer close(resultChannel)
	// 	workerControl.TaskUnparsor(r)

	// 	newTask := &WorkerCont.Task{
	// 		FunctionName: WorkerCont.GetStatus,
	// 		TaskSpecific: workerControl,
	// 	}

	// 	taskPool.WorkerAllocate(newTask)
	// 	result := <-resultChannel
	// 	encoder := json.NewEncoder(w)
	// 	encoder.Encode(result)
	// })

	http.HandleFunc("/CreateVM", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // POST로 요청 제한
			http.Error(w, "Control : Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		param, err := util.UnmarshalBodyAndClose[WorkerCont.CreateVMParam](r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400 오류
			println("Control : Failed to read or parse JSON")
			json.NewEncoder(w).Encode(WorkerCont.ControlError{
				Message: "Failed to create VM",
				Errors:  "Control : Failed to read or parse JSON",
			})
			return
		}
		_, err = contextStruct.AssignInternalAddress()
		param.Network.Ips = []string{"10.5.15.10"}
		fmt.Println(param.Network.Ips)
		excludeFields := map[string]bool{"Network": true}
		//var privateKey string
		for i := range param.Users {
			privateKey, publicKey, err := WorkerCont.SshKeygen()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				println("Control : ssh err")
				json.NewEncoder(w).Encode(WorkerCont.ControlError{
					Message: "Failed to create VM",
					Errors:  "Control : ssh err",
				})
				return
			}
			if len(param.Users[i].Ssh) == 0 {
				param.Users[i].Ssh = make([]string, 1)
				param.Users[i].Ssh[0] = publicKey
			} else {
				param.Users[i].Ssh = append(param.Users[i].Ssh, publicKey)
			}
			//WorkerCont.SshPrivateStore(param.UUID+"_"+strconv.Itoa(i), privateKey)
			WorkerCont.GuacamoleConfig(param.Users[0].UserName, param.UUID, param.Network.Ips[0], privateKey)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(WorkerCont.ControlError{
				Message: "Failed to create VM",
				Errors:  "Control : Failed to assign internal address",
			})
			return
		}

		err = WorkerCont.ValidateStruct(param, excludeFields)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			println("Control : Invalid parameters provided")
			json.NewEncoder(w).Encode(WorkerCont.ControlError{
				Message: "Failed to create VM",
				Errors:  "Control : Invalid parameters provided",
			})
			return
		}

		core := WorkerCont.UpdateCoreAndSelectCoreForNewVM(contextStruct, uint64(param.HWInfo.Memory), float64(param.HWInfo.CPU))
		task := WorkerCont.NewCreateVMTask(core, param)
		resp, err := task.Await()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			println("Core : Forced error response instead of core response")
			json.NewEncoder(w).Encode(WorkerCont.ControlError{
				Message: "Failed to create VM",
				Errors:  err.Error(),
			})
			return
		}

		if resp.Errors.ErrorType != "" {
			println("Core : Forced error response instead of core response")
			WorkerCont.Errorhandler(resp)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(WorkerCont.ControlError{
				Message: "Failed to create VM",
				//Errors:  "Core : Forced error response instead of core response(%s)", resp,
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(WorkerCont.ControlError{
			Message: "Success to create VM",
			Errors:  "",
		})
	})

	// http.HandleFunc("/DeleteVM", func(w http.ResponseWriter, r *http.Request) {
	// 	param, err := util.UnmarshalBodyAndClose[WorkerCont.DeletevmParam](r.Body)
	// 	fmt.Printf("%v\n", param)
	// 	if err != nil {
	// 		http.Error(w, "Failed to read or parse JSON", http.StatusBadRequest)
	// 		return
	// 	}
	// 	task := WorkerCont.NewDeleteVMTask(&vms.Core{IP: "223.194.20.119", Port: 28779}, param)
	// 	resp, err := task.Await()
	// 	if err != nil {
	// 		http.Error(w, "Failed to Delete VM", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	//task.errorhandler()

	// 	encoder := json.NewEncoder(w)
	// 	if err = encoder.Encode(resp); err != nil {
	// 		http.Error(w, "Failed to encode result", http.StatusInternalServerError)
	// 		return
	// 	}

	// })
	// http.HandleFunc("GET /ConnectVM", func(w http.ResponseWriter, b *http.Request) {
	// 	taskPool.WorkerAllocate(&WorkerCont.Task{
	// 		FunctionName: WorkerCont.ConnectV,
	// 	})
	// })
	// http.HandleFunc("GET /CheckVMHealth", func(w http.ResponseWriter, b *http.Request) {
	// 	taskPool.WorkerAllocate(&WorkerCont.Task{
	// 		FunctionName: WorkerCont.UpdateStat,
	// 	})
	// })

	err := http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
	if err != nil {
		return err
	}

	return nil
}
