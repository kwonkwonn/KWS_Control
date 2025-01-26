package WorkerCont

import (
	//"context"
	"fmt"
	_ "log"
	_ "sync"
	"time"
)

func PseudoRequestSender(workerHandler *TaskHandler) {
	//workerHandler.Lock()
	con := Task{
		FunctionName: ConnectV,
	}
	del := Task{
		FunctionName: DeleteV,
	}
	cre := Task{
		FunctionName: CreateV,
	}
	up := Task{
		FunctionName: UpdateStat,
	}
	stat := Task{
		FunctionName: GetStatus,
	}
	taskList := []Task{con, del, cre, up, stat}
	fmt.Print(len(taskList))
	for j := 0; j < 20; j++ {
		workerHandler.WorkerAllocate(&taskList[j%(len(taskList))])
	}
	//workerHandler.Unlock()
}

func (t *TaskWorker) UpdateStatusTest() {
	t.workDescription(UpdateStat)
	time.Sleep(5 * time.Second)
}
func (t *TaskWorker) CreateVMTest() {
	t.workDescription(CreateV)
	time.Sleep(5 * time.Second)
}
func (t *TaskWorker) ConnectVMTest() {
	t.workDescription(ConnectV)
	time.Sleep(5 * time.Second)
}
func (t *TaskWorker) DeleteVMTest() {
	t.workDescription(DeleteV)
	time.Sleep(5 * time.Second)
}
func (t *TaskWorker) GetStatusTest() {
	t.workDescription(DeleteV)
	time.Sleep(5 * time.Second)
}

func (w *TaskWorker) workDescription(workName functionName) {
	fmt.Printf("\n")
	fmt.Println("*******************************************************")
	fmt.Printf("worker %d \n", w.workerNum)
	fmt.Printf("current length of workLoad is %d currently working on %s \n", w.tasksLength, functionNameEmmitor(workName))
	fmt.Println("*******************************************************")
	fmt.Printf("\n")
}

func functionNameEmmitor(functionNum functionName) string {
	switch functionNum {
	case CreateV:
		return "CreateVM"
	case UpdateStat:
		return "Update status"
	case ConnectV:
		return "connect VM"
	case DeleteV:
		return "delete VM"
	case GetStatus:
		return "status of VM"
	default:
		return "undefined Task"
	}
}
