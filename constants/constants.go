package constants

const (
	TRACE_ID_HTTP_HEADER           = "X-B3-Traceid"
	LOGGER_KEY                     = "logger"
	TRACE_KEY                      = "traceID"
	JSON_FORMATTER_HTTP_HEADER     = "JSON"
	NO_TRACE_ID                    = "no-trace-id"
	JSON                           = "json"
	TRACING_SESSION_HEADER_KEY     = "Session-Tracing-ID"
	TRACE_CONFIG_MAX_ANNOTATIONS   = 128

	TRACING_CLIENT_PUBLIC_IP_HEADER = "X-Original-Forwarded-For"
	TRACING_CLIENT_PUBLIC_IP        = "client-public-ip"
	TRACING_SESSION_ID              = "session_tracing_id"
	TRACING_APP_VERSION_HEADER_KEY  = "App-Version"
	TRACING_APP_VERSION             = "app_version"
	TRACING_DEVICE_INFO_HEADER_KEY  = "Device-Info"
	TRACING_DEVICE_INFO             = "device_info"
)
