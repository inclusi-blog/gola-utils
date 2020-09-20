package mask_util

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MaskUtilsTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestMaskUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(MaskUtilsTestSuite))
}

func (suite *MaskUtilsTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "dummyUrl", nil)
}

func (suite *MaskUtilsTestSuite) TestMaskEmail_ShouldReturnMaskedEmail() {
	expected := "a*c@g***l.com"
	email := "abc@gmail.com"

	actual := MaskEmail(suite.context, email)

	suite.Equal(expected, actual)

	expected = "g******h@y***o.co.in"
	email = "gobinath@yahoo.co.in"

	actual = MaskEmail(suite.context, email)

	suite.Equal(expected, actual)
}

func (suite *MaskUtilsTestSuite) TestMaskEmail_ShouldReturnErrorIfEmailValueIsEmpty() {
	email := ""

	actual := MaskEmail(suite.context, email)

	suite.Equal(fallbackText, actual)
}
