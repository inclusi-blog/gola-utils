package http

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/url"
	"testing"
)

type IntrospectionRequestBuilderTest struct {
	suite.Suite
	introspectionRequestBuilder IntrospectionRequestBuilder
}

func TestIntrospectionRequestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(IntrospectionRequestBuilderTest))
}

func (suite *IntrospectionRequestBuilderTest) SetupTest() {
	suite.introspectionRequestBuilder = NewIntrospectionRequestBuilder()
}

func (suite IntrospectionRequestBuilderTest) TestIntrospectionRequestBuilderShouldBuildHttpRequestForIntrospection() {
	request, error := suite.introspectionRequestBuilder.Build("admin_url", "my_tok##@en")
	suite.Nil(error)
	suite.Equal("application/x-www-form-urlencoded", request.Header.Get("Content-Type"))
	suite.Equal("application/json", request.Header.Get("Accept"))
	suite.Equal("admin_url/oauth2/introspect", request.URL.Path)
	suite.Equal("POST", request.Method)
	bytes, error := ioutil.ReadAll(request.Body)
	expectedBodyString := fmt.Sprintf("token=%s", url.QueryEscape("my_tok##@en"))
	suite.Equal([]byte(expectedBodyString), bytes)
}
