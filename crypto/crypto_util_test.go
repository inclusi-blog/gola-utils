package crypto

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CryptoUtilTestSuite struct {
	suite.Suite
	recorder   *httptest.ResponseRecorder
	context    *gin.Context
	cryptoUtil CryptoUtil
}

func TestCryptoUtilTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoUtilTestSuite))
}

func (suite *CryptoUtilTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "dummyUrl", nil)
}

func (suite CryptoUtilTestSuite) TestDecipherShouldReturnErrorWhenEncryptedTextIsEmpty() {
	suite.cryptoUtil = NewCryptoUtil("crypto-svc-url")
	_, err := suite.cryptoUtil.Decipher(suite.context, "")
	suite.Equal(errors.New("text is empty"), err)
}

func (suite CryptoUtilTestSuite) TestDecipherShouldReturnErrorWhenCryptoServiceReturnsErrorForDecryptText() {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.cryptoUtil = NewCryptoUtil(testServer.URL)
	_, err := suite.cryptoUtil.Decipher(suite.context, "encrypted_text")
	suite.NotNil(err)
}

func (suite CryptoUtilTestSuite) TestDecipherShouldReturnErrorWhenCryptoServiceReturnsInvalidResponse() {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.cryptoUtil = NewCryptoUtil(testServer.URL)
	_, err := suite.cryptoUtil.Decipher(suite.context, "encrypted_text")
	suite.NotNil(err)
}

func (suite CryptoUtilTestSuite) TestDecipherShouldReturnDecryptedText() {
	body, _ := json.Marshal(model.CryptoResponse{DecryptedText: "decrypted_text"})
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write(body)
	}))
	defer func() { testServer.Close() }()
	suite.cryptoUtil = NewCryptoUtil(testServer.URL)
	actualText, err := suite.cryptoUtil.Decipher(suite.context, "encrypted_text")
	suite.Equal("decrypted_text", actualText)
	suite.Nil(err)
}
