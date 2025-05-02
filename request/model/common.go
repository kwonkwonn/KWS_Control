package model

type CoreResponse[T any] struct {
	Information *T              `json:"information,omitempty"`
	Message     string          `json:"message"`
	Errors      ErrorDescriptor `json:"errors,omitempty"`
}

type ErrorDescriptor struct {
	ErrorType VirError `json:"error type"`
	Detail    error    `json:"detail"`
}

type VirError string

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
