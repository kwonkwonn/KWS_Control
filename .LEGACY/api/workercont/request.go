package WorkerCont

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
)

type CoreRequestTask[P any, R any] struct {
	Core     *vms.Core
	Method   string
	Endpoint string
	Request  P
	Response R
}

func (t *CoreRequestTask[P, R]) Await() (body R, err error) {
	jsonData, err := json.Marshal(t.Request)
	if err != nil {
		println("Struct encoding error")
		return
	}
	//fmt.Println(t.Request)
	requestUrl := url.URL{
		Scheme: "http",
		Host:   t.Core.IP + ":" + strconv.Itoa(t.Core.Port),
		Path:   t.Endpoint,
	}
	fmt.Println(requestUrl.String())
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(t.Method, requestUrl.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		println("Control : Failed to create request")
		//
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 요청 실행
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Control : Core access error - Server did not respond in time")
		err = fmt.Errorf("server timeout or connection error")
		return
	}
	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if err == nil {
			err = e
		}
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("Core Response:", string(b))
	json.Unmarshal(b, &body)
	return
}

func NewCreateVMTask(core *vms.Core, param CreateVMParam) CoreRequestTask[CreateVMParam, CoreResponse[CreateVMResp]] {
	return CoreRequestTask[CreateVMParam, CoreResponse[CreateVMResp]]{
		Core:     core,
		Method:   "POST",
		Endpoint: "/createVM",
		Request:  param,
	}
}

func NewDeleteVMTask(core *vms.Core, param DeletevmParam) CoreRequestTask[DeletevmParam, CoreResponse[DeleteVMResp]] {
	return CoreRequestTask[DeletevmParam, CoreResponse[DeleteVMResp]]{
		Core:     core,
		Method:   "POST",
		Endpoint: "/DeleteVM",
		Request:  param,
	}
}

func NewGetCoreMachineFreeCpuInfoTask(core *vms.Core) CoreRequestTask[GetMachineStatusParam, CoreResponse[CoreMachineCpuInfoResp]] {
	return CoreRequestTask[GetMachineStatusParam, CoreResponse[CoreMachineCpuInfoResp]]{
		Core:     core,
		Method:   "GET",
		Endpoint: "/getStatusHost",
		Request:  GetMachineStatusParam{HostDataType: CpuInfo},
	}
}

func NewGetCoreMachineFreeMemoryInfoTask(core *vms.Core) CoreRequestTask[GetMachineStatusParam, CoreResponse[CoreMachineMemoryInfoResp]] {
	return CoreRequestTask[GetMachineStatusParam, CoreResponse[CoreMachineMemoryInfoResp]]{
		Core:     core,
		Method:   "GET",
		Endpoint: "/getStatusHost",
		Request:  GetMachineStatusParam{HostDataType: MemInfo},
	}
}
