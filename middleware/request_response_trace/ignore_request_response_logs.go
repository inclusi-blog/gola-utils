package request_response_trace

type IgnoreRequestResponseLogs struct {
	PartialApiPath       string
	IsRequestLogAllowed  bool
	IsResponseLogAllowed bool
}
