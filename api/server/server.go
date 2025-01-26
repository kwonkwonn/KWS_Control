package server

import (
	"encoding/json"
	"fmt"

	//"io"
	"net/http"
	"strconv"

	WorkerCont "github.com/easy-cloud-Knet/KWS_Control/api/workercont"
	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func Server(portNum int, taskPool *WorkerCont.TaskHandler, contextStruct *vms.ControlInfra) error {
	// main server와 통신하기 위한 http 서버
	// gin.DefaultWriter = io.Discard
	http.HandleFunc("Get /getStatus", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		workerControl := &WorkerCont.TaskControlGetStatus{
			ResultChan: make(chan WorkerCont.TaskExecutionResult),
		}
		resultChannel := workerControl.ResultChan
		defer close(resultChannel)
		workerControl.TaskUnparsor(r)

		newTask := &WorkerCont.Task{
			FunctionName: WorkerCont.GetStatus,
			TaskSpecific: workerControl,
		}

		taskPool.WorkerAllocate(newTask)
		result := <-resultChannel
		encoder := json.NewEncoder(w)
		encoder.Encode(result)
	})

	http.HandleFunc("/CreateVM", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { // POST로 요청 제한
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// TaskControlCreateVM 생성 및 요청 파싱
		workerControl := &WorkerCont.TaskControlCreateVM{
			ResultChan: make(chan string),
			Vms:        contextStruct,
		}
		resultChannel := workerControl.ResultChan
		defer close(resultChannel)

		// JSON 파싱 및 에러 처리
		if err := workerControl.TaskUnparsor(r); err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			fmt.Printf("Error in TaskUnparsor: %v\n", err) // 필요 시 유지
			return
		}

		// Task 생성 및 작업 할당
		newTask := &WorkerCont.Task{
			FunctionName: WorkerCont.CreateV,
			TaskSpecific: workerControl,
		}
		taskPool.WorkerAllocate(newTask)

		// 결과 처리 및 응답
		result := <-resultChannel
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(result); err != nil {
			http.Error(w, "Failed to encode result", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/DeleteVM", func(w http.ResponseWriter, r *http.Request) {
		workerControl := &WorkerCont.TaskControlDeleteVM{
			ResultChan: make(chan string),
			//Vms:        contextStruct,
		}
		resultChannel := workerControl.ResultChan
		defer close(resultChannel)

		// JSON 파싱 및 에러 처리
		if err := workerControl.TaskUnparsor(r); err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			fmt.Printf("Error in TaskUnparsor: %v\n", err) // 필요 시 유지
			return
		}

		// Task 생성 및 작업 할당
		newTask := &WorkerCont.Task{
			FunctionName: WorkerCont.DeleteV,
			TaskSpecific: workerControl,
		}
		taskPool.WorkerAllocate(newTask)

		// 결과 처리 및 응답
		result := <-resultChannel
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(result); err != nil {
			http.Error(w, "Failed to encode result", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("GET /ConnectVM", func(w http.ResponseWriter, b *http.Request) {
		taskPool.WorkerAllocate(&WorkerCont.Task{
			FunctionName: WorkerCont.ConnectV,
		})
	})
	http.HandleFunc("GET /CheckVMHealth", func(w http.ResponseWriter, b *http.Request) {
		taskPool.WorkerAllocate(&WorkerCont.Task{
			FunctionName: WorkerCont.UpdateStat,
		})
	})

	err := http.ListenAndServe(":"+strconv.Itoa(portNum), nil)
	if err != nil {
		return err
	}

	return nil
}

/*
회의 해봐야 하는 내용들
1. VM 생성 완료 했을 때 벡에다가 리턴해야 하는게 뭔지?
2. Core 컴퓨터가 실행되면 Control에 Core 정보 보내줘야함.

*/
