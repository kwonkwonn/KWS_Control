package WorkerCont

import (
	"context"
	"fmt"
	"time"
)

//각각의 worker노드와 통신하기 위한 클라이언트 서버
type functionName int32

const NumOfTaskHandler = 5
const NumOfAllWorkload = 10

const (
	InitWorker functionName = iota + 1
	UpdateStat
	CreateV
	ConnectV
	DeleteV
	GetStatus
)

func InitWorkers(pool *TaskHandler) {
	//pool *TaskHandler as args
	pool.TaskHandlersList = make([]*TaskWorker, NumOfTaskHandler)
	pool.workingIndex = 0
	// pool.TaskPool=make([]*Task,10)

	for i := 0; i < NumOfTaskHandler; i++ {
		pool.TaskHandlersList[i] = &TaskWorker{
			tasksLength: 0,
			workLoads:   make(chan *Task, NumOfAllWorkload),
			//스레드 관련 내용인거 같은데 뭔소리인지 모르겠음.
			//내생각엔 스레드를 처리하는 큐? 느낌인거 같음.


			workerNum:   i, // 코어의 고유 번호
		}
		go pool.TaskHandlersList[i].StartWorking()//????

	}
}

func (t *TaskWorker) StartWorking() {
	for {
		ctx := context.Background()
		select {
		case work, ok := <-t.workLoads:
			{
				if !ok {
					fmt.Println("channel closed")
					return
				} else {
					switch work.FunctionName {
					case CreateV:
						t.CreateVMTest(ctx)
						t.tasksLength--
					case UpdateStat:
						t.UpdateStatusTest(ctx)
						t.tasksLength--
					case ConnectV:
						t.ConnectVMTest(ctx)
						t.tasksLength--
					case DeleteV:
						t.DeleteVMTest(ctx)
						t.tasksLength--
					case GetStatus:
						t.GetStatus(ctx, work)
						t.tasksLength--
					default:
						fmt.Printf("undefined task")
					}
				}
			}
		default:
			// fmt.Println("work Done, waiting")
			time.Sleep(time.Microsecond * 300)
		}
	}
}

func (t *TaskHandler) WorkerAllocate(task *Task) {
	for {
		workerIndex := t.workingIndex % NumOfTaskHandler
		worker := t.TaskHandlersList[workerIndex]

		worker.taskLenMu.Lock()
		if worker.tasksLength < 9 {
			worker.workLoads <- task
			worker.tasksLength++
			fmt.Printf("Allocated task to worker %d. Current tasks: %d\n", workerIndex, worker.tasksLength)
			worker.taskLenMu.Unlock()
			t.workingIndex = (t.workingIndex + 1) % NumOfTaskHandler
			return
		}
		worker.taskLenMu.Unlock()
		t.workingIndex = (t.workingIndex + 1) % NumOfTaskHandler
	}
}
