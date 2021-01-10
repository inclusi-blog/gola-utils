package error

import "net/http"

type AuthenticationError struct {
}

func (e AuthenticationError) Error() string {
	return "Invalid access token / not present"
}

type HydraInternalServerError struct {
}

func (e HydraInternalServerError) Error() string {
	return "Internal error in Hydra"
}

type additionalData interface{}

type Error struct {
	ErrorCode      string         `json:"errorCode,omitempty"`
	ErrorMessage   string         `json:"errorMessage"`
	AdditionalData additionalData `json:"additionalData,omitempty"`
}

type OAuthMiddlewareError struct {
	HttpStatusCode int
	ErrorResponse  Error
}

func (r Error) Error() string {
	return "ErrorCode: " + r.ErrorCode +
		"ErrorMessage: " + r.ErrorMessage
}

const (
	InvalidAccessTokenErrorCode = "ERR_INVALID_ACCESS_TOKEN"
	InvalidIdTokenErrorCode     = "ERR_INVALID_ID_TOKEN_ERROR"
	InternalServerErrorCode     = "ERR_INTERNAL_SERVER_ERROR_CODE"
)

func InternalServerErrorFunc(errorMessage string) *OAuthMiddlewareError {
	return &OAuthMiddlewareError{HttpStatusCode: http.StatusInternalServerError,
		ErrorResponse: Error{ErrorCode: InternalServerErrorCode, ErrorMessage: errorMessage}}
}

var (
	InvalidIdTokenError = &OAuthMiddlewareError{
		HttpStatusCode: http.StatusUnauthorized,
		ErrorResponse:  Error{ErrorCode: InvalidIdTokenErrorCode, ErrorMessage: "Id token invalid"},
	}
	InvalidAccessTokenError = &OAuthMiddlewareError{
		HttpStatusCode: http.StatusUnauthorized,
		ErrorResponse:  Error{ErrorCode: InvalidAccessTokenErrorCode, ErrorMessage: "Invalid access token / not present"},
	}
)
