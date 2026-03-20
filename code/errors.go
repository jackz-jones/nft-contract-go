package code

import (
	"encoding/json"
	"fmt"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 一些内部错误信息修饰符字符串
const (
	InternalDataStoreFailedStrModifier     = "failed to putstate[%s %s],err: %v"
	InternalDataReadFailedStrModifier      = "failed to getstate[%s %s],err: %v"
	InternalAtoiFailedStrModifier          = "failed to convert %s(%v) to int,err: %v"
	InternalJsonMarshalFailedStrModifier   = "failed to json marshal %v,err:%v"
	InternalJsonUnmarshalFailedStrModifier = "failed to json unmarshal %s,err:%v"
	InternalEmitEventFailedStrModifier     = "failed to emit %s,err: %v"
)

type ErrCode int

const (
	ErrGetOriginFailed ErrCode = iota + 720000
	ErrInternalDataStoreFailed
	ErrInternalDataReadFailed
	ErrInternalAtoiFailed
	ErrInternalJsonMarshalFailed
	ErrInternalJsonUnmarshalFailed
	ErrInvalidMethod
	ErrInternalGetSenderFailed
	ErrInternalEmitEvent
	ErrInvalidParam
	ErrNFTLocked
	ErrUnknownState
)

var errorMessages = map[ErrCode]string{
	ErrGetOriginFailed:             "get origin failed",
	ErrInternalDataStoreFailed:     "internal data store failed",
	ErrInternalDataReadFailed:      "internal data read failed",
	ErrInternalAtoiFailed:          "internal atoi failed",
	ErrInternalJsonMarshalFailed:   "internal json marshal failed",
	ErrInternalJsonUnmarshalFailed: "internal json unmarshal failed",
	ErrInvalidMethod:               "invalid method",
	ErrInternalGetSenderFailed:     "internal get sender failed",
	ErrInternalEmitEvent:           "internal emit event failed",
	ErrInvalidParam:                "invalid param",
	ErrNFTLocked:                   "nft locked",
	ErrUnknownState:                "unknown state",
}

func (ec ErrCode) String() string {
	return errorMessages[ec]
}

// Error represents a standard code structure for the smart contract.
// It includes an error code, a message, and detailed information about the code.
type Error struct {
	// The error code corresponding to the type of code encountered.
	ErrorCode int `json:"errorCode"`
	// A message describing the code.
	ErrorMessage string `json:"errorMessage"`
	// Detailed information about the code, which may include specifics about the cause.
	ErrorDetail string `json:"errorDetail"`
}

// NewError creates a new Error instance with the given error code and detail.
func NewError(errorCode ErrCode, errorDetail string) *Error {
	return &Error{
		ErrorCode:    int(errorCode),
		ErrorMessage: errorMessages[errorCode],
		ErrorDetail:  errorDetail,
	}
}

// ErrorResp creates a new Error instance with the given error code and detail.
func ErrorResp(errorCode ErrCode, errorDetail string) protogo.Response {
	errStr, err := json.Marshal(NewError(errorCode, errorDetail))
	if err != nil {
		return sdk.Error(fmt.Sprintf("failed to marshal code, code:%d, detail:%s, %v", errorCode, errorDetail, err))
	}
	return sdk.Error(string(errStr))
}
