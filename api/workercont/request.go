package WorkerCont

import (
	"bytes"
	"encoding/json"
	vms "github.com/easy-cloud-Knet/KWS_Control/vm"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type CoreRequestTask[P any, R any] struct {
	Core     *vms.Core
	Endpoint string
	Request  P
	Response R
}

func (t *CoreRequestTask[P, R]) Await() (body R, err error) {
	jsonData, err := json.Marshal(t.Request)
	if err != nil {
		return
	}

	requestUrl := url.URL{
		Scheme: "http",
		Host:   t.Core.IP + ":" + strconv.Itoa(t.Core.Port),
		Path:   t.Endpoint,
	}

	resp, err := http.Post(requestUrl.String(), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
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

	err = json.Unmarshal(b, &body)

	return
}

func NewCreateVMTask(core *vms.Core, param CreateVMParam) CoreRequestTask[CreateVMParam, string] {
	return CoreRequestTask[CreateVMParam, string]{
		Core:     core,
		Endpoint: "/createVM",
		Request:  param,
	}
}
