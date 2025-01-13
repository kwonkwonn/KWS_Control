package WorkerCont

import (
	"encoding/json"
	"net/http"
	"sync"

	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

type GetStatusParam struct {
	UUID string `json:"UUID"`
}

type CreateVMParam struct {
	UUID string `json:"UUID"`
	RAM  int    `json:"RAM"`
	CPU  int    `json:"CPU"`
}

//parameters for each functions

type Task struct {
	FunctionName functionName
	TaskSpecific TaskJustifier
}
type TaskWorker struct {
	taskLenMu   sync.Mutex
	tasksLength int
	workLoads   chan *Task
	workerNum   int
}

type TaskHandler struct {
	TaskHandlersList []*TaskWorker // 현재 돌아가고 있는 스레드를 가리키는 포인터
	workingIndex     int // 현재 돌아가고 있는 스레드의 갯수
} //테스크 임베딩, 필드 추가 필요

type TaskControlCreateVM struct {
	ResultChann chan TaskExecutionResult
	UUID        vms.UUID
	RAM         int
	CPU         int
	//추가 필요
}

type TaskControlGetStatus struct {
	ResultChan chan TaskExecutionResult
	UUID       vms.UUID
}

type TaskJustifier interface {
	TaskUnparsor(r *http.Request) error
}

//functions for each structures, needed for interface
func (t *TaskControlCreateVM) TaskUnparsor(r *http.Request) error {
	var param CreateVMParam
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		return err
	}
	t.RAM = param.RAM
	t.CPU = param.CPU
	return nil
}

func (t *TaskControlGetStatus) TaskUnparsor(r *http.Request) error {
	var param GetStatusParam
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		//에러 정의 필요
		return err
	}
	return nil
}

type TaskExecutionResult struct {
	IsSuccess   bool
	InVMContext vms.VMInfo
}
