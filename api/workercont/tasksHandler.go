package WorkerCont

import (
	//"context"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	//vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

func (t *TaskWorker) UpdateStatus() {

}
func (t *TaskWorker) CreateVM(CreateVM *TaskControlCreateVM) {
	fmt.Printf("VM Name: %s\n", CreateVM.Param.DomName)

	// JSON 변환
	jsonData, err := json.Marshal(CreateVM.Param)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	var apiURL = "http://223.194.20.119:28779/createVM"
	// HTTP POST 요청 생성
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		CreateVM.ResultChan <- fmt.Sprintf("Error sending request to API: %v", err)
		return
	}
	defer resp.Body.Close() // 응답 본문 닫기

	// 응답 확인
	if resp.StatusCode == http.StatusOK {
		CreateVM.ResultChan <- "VM successfully created and assigned!"
	} else {
		CreateVM.ResultChan <- fmt.Sprintf("Failed to create VM. Status code: %d", resp.StatusCode)
	}
}
func (t *TaskWorker) ConnectVM() {
}
func (t *TaskWorker) DeleteVM(DeleteVM *TaskControlDeleteVM) {
	jsonData, err := json.Marshal(DeleteVM.Param)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println("1111")
	var apiURL = "http://223.194.20.119:28779/DeleteVM"
	// HTTP POST 요청 생성
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		DeleteVM.ResultChan <- fmt.Sprintf("Error sending request to API: %v", err)
		return
	}
	defer resp.Body.Close() // 응답 본문 닫기

	// 응답 확인
	if resp.StatusCode == http.StatusOK {
		DeleteVM.ResultChan <- "VM successfully created and assigned!"

	} else {
		DeleteVM.ResultChan <- fmt.Sprintf("Failed to create VM. Status code: %d", resp.StatusCode)
	}
}
func (t *TaskWorker) GetStatus(task *Task) {

}
