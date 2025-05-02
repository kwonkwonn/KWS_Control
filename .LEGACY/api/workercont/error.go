package WorkerCont

import (
	"fmt"
)

type VirError string
type ConError string

const (
	Unparse ConError = "parse Error"
)

const (
	FaildDeEncoding   VirError = "Error Not Found"
	DomainSearchError VirError = "Error Searching Domain"
	NoSuchDomain      VirError = "Domain Not Found"

	DomainGenerationError VirError = "error Generating Domain" // 이 부분 추가
	LackCapacityRAM       VirError = "Not enough RAM"          // control
	LackCapacityCPU       VirError = "Not Enough CPU"          //
	LackCapacityHD        VirError = "Not Enough HardDisk"     //

	InvalidUUID VirError = "Invalid UUID Provided"

	InvalidParameter VirError = "Invalid parameter entered"
	WrongParameter   VirError = "Not validated parameter In"

	DomainStatusError VirError = "Error Retrieving Domain Status"
	HostStatusError   VirError = "Error Retrieving Host Status"

	DeletionDomainError VirError = "Error Deleting Domain"
	DomainShutdownError VirError = "Failed in Shutting down domain"
)

type ControlError struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type CoreResponse[T any] struct {
	Information *T              `json:"information,omitempty"`
	Message     string          `json:"message"`
	Errors      ErrorDescriptor `json:"errors,omitempty"`
}

type ErrorDescriptor struct {
	ErrorType VirError `json:"error type"`
	Detail    error    `json:"detail"`
}

type CreateVMResp struct {
	State     string `json:"state"`
	MaxMem    uint64 `json:"maxmem"`
	Memory    uint64 `json:"memory"`
	NrVirtCpu uint   `json:"nrVirtCpu"`
	CpuTime   uint64 `json:"cpuTime"`
}

type DeleteVMResp struct {
	State     string `json:"state"`
	MaxMem    uint64 `json:"maxmem"`
	Memory    uint64 `json:"memory"`
	NrVirtCpu uint   `json:"nrVirtCpu"`
	CpuTime   uint64 `json:"cpuTime"`
}

func Errorhandler[T any](Errorcode CoreResponse[T]) {
	// 에러가 존재하지 않으면 함수 종료
	if Errorcode.Errors.ErrorType == "" {
		fmt.Println("No errors detected.")
		return
	}

	// VirError 종류에 따른 처리
	switch Errorcode.Errors.ErrorType {
	case FaildDeEncoding:
		fmt.Println("[ERROR] Resource not found.")
	case DomainSearchError:
		fmt.Println("[ERROR] Domain search failed.")
	case NoSuchDomain:
		fmt.Println("[ERROR] Specified domain does not exist.")
	case DomainGenerationError: // 여기 추가
		fmt.Println("[ERROR] Error generating domain.")
	case LackCapacityRAM:
		fmt.Println("[ERROR] Insufficient RAM capacity.")
	case LackCapacityCPU:
		fmt.Println("[ERROR] Insufficient CPU capacity.")
	case LackCapacityHD:
		fmt.Println("[ERROR] Insufficient Hard Disk capacity.")
	case InvalidUUID:
		fmt.Println("[ERROR] Invalid UUID provided.")
	case WrongParameter:
		fmt.Println("[ERROR] Invalid parameter received.")
	case DomainStatusError:
		fmt.Println("[ERROR] Failed to retrieve domain status.")
	case DeletionDomainError:
		fmt.Println("[ERROR] Failed to delete domain.")
	case DomainShutdownError:
		fmt.Println("[ERROR] Failed to shut down domain.")
	default:
		fmt.Println("[ERROR] Unknown error occurred.")
	}

	// 상세한 에러 정보가 있는 경우 출력
	if Errorcode.Errors.Detail != nil {
		fmt.Println("Detailed Error:", Errorcode.Errors.Detail)
	}
}
