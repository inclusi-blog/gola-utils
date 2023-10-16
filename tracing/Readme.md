## Tracing

#### tools 
1. OpenCensus exporter
2. OpenCensus collector
3. Jaeger backend

#### Steps to add tracing


1. Initialise a exporter to export traces
```go 
import "github.com/inclusi-blog/gola-utils/http/request"

tracing.Init("service-name", "collector-address")
```
2. Instrument server
```go
app := tracing.WithTracing(router, "api-healthz-endpoint") // to avoid tracing health endpoint
gohttp.ListenAndServe(":8080", app)
```
3. Instrument client for inter service calls
```go
import "go.opencensus.io/plugin/ochttp"
httpClient := &http.Client{Transport: &ochttp.Transport{}}
```
