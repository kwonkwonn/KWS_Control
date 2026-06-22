package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

// ApiCreateVmRequest for POST /vm HTTP Request Body contract.
type ApiCreateVmRequest struct {
	DomType    string          `json:"domType"`
	DomName    string          `json:"domName"`
	UUID       structure.UUID  `json:"uuid"`
	OS         string          `json:"os"`
	HWInfo     ApiHardwareInfo `json:"HWInfo"`
	Network    ApiNetworkInfo  `json:"network"`
	Users      []ApiUserInfo   `json:"users"`
	SubnetType string          `json:"Subnettype"`
}

type ApiHardwareInfo struct {
	CPU    uint32 `json:"cpu"`
	Memory uint32 `json:"memory"` // MiB
	Disk   uint32 `json:"disk"`   // MiB
}

type ApiNetworkInfo struct {
	IPs []string `json:"ips"`
	// NetType은 내부에서 0 고정 — API 클라이언트가 전송하더라도 무시됨
}

type ApiUserInfo struct {
	Name              string   `json:"name"`
	Groups            string   `json:"groups"`
	Password          string   `json:"passWord"`
	SSHAuthorizedKeys []string `json:"ssh"`
}

// ToServiceInput은 HTTP DTO를 서비스 계층 DTO로 변환
func (r *ApiCreateVmRequest) ToServiceInput() service.CreateVMInput {
	users := make([]service.UserSpec, len(r.Users))
	for i, u := range r.Users {
		users[i] = service.UserSpec{
			Name:              u.Name,
			Groups:            u.Groups,
			Password:          u.Password,
			SSHAuthorizedKeys: u.SSHAuthorizedKeys,
		}
	}
	return service.CreateVMInput{
		UUID:    r.UUID,
		DomType: r.DomType,
		DomName: r.DomName,
		OS:      r.OS,
		HardwareInfo: service.HardwareSpec{
			CPU:    r.HWInfo.CPU,
			Memory: r.HWInfo.Memory,
			Disk:   r.HWInfo.Disk,
		},
		Network: service.NetworkSpec{
			IPs: r.Network.IPs,
		},
		Users:      users,
		SubnetType: r.SubnetType,
	}
}

func (c *handlerContext) createVm(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req ApiCreateVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("createVm: failed to parse request body: %v", err, true)
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			util.RespondError(w, http.StatusBadRequest, "invalid JSON format in request body")
		} else {
			util.RespondError(w, http.StatusBadRequest, "err req body parsing: "+err.Error())
		}
		return
	}

	if req.HWInfo.Memory == 0 || req.HWInfo.CPU == 0 || req.HWInfo.Disk == 0 {
		util.RespondError(w, http.StatusBadRequest, "Memory, CPU, and Disk must be non-zero")
		return
	}

	if err := service.CreateVM(req.ToServiceInput(), c.context, c.rdb); err != nil {
		log.Error("createVm: failed to create VM: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusCreated, map[string]string{"message": "VM created successfully"})
}
