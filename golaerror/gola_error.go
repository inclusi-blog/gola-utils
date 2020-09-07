package golaerror

type additionalData interface{}

type Error struct {
	ErrorCode      string         `json:"errorCode"`
	ErrorMessage   string         `json:"errorMessage"`
	AdditionalData additionalData `json:"additionalData,omitempty"`
}

func (r Error) Error() string {
	return "ErrorCode: " + r.ErrorCode + " ErrorMessage: " + r.ErrorMessage
}

func New(errorCode, errorMessage string, additionalData interface{}) Error {
	return Error{
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
		AdditionalData: additionalData,
	}
}
