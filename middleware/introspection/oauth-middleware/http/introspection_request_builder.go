package http

//go:generate mockgen -source=introspection_request_builder.go -destination=./../mocks/mock_introspection_request_builder.go -package=mocks

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IntrospectionRequestBuilder interface {
	Build(hydraAdminBaseUrl string, accessToken string) (*http.Request, error)
}

type introspectionRequestBuilder struct {
}

func NewIntrospectionRequestBuilder() IntrospectionRequestBuilder {
	return introspectionRequestBuilder{}
}

func (introspectionRequestBuilder introspectionRequestBuilder) Build(hydraAdminBaseUrl string, accessToken string) (*http.Request, error) {
	headers := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"application/json"},
	}

	data := url.Values{}
	data.Set("token", accessToken)

	introspectionUrl := fmt.Sprintf("%s/oauth2/introspect", hydraAdminBaseUrl)
	request, requestError := http.NewRequest("POST", introspectionUrl, strings.NewReader(data.Encode()))
	if requestError != nil {
		return nil, requestError
	}

	request.Header = headers
	return request, nil
}
