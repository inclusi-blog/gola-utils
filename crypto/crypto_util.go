package crypto

import (
	"context"
	"errors"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/inclusi-blog/gola-utils/http/request"
	"github.com/inclusi-blog/gola-utils/http/util"
	"github.com/inclusi-blog/gola-utils/model"
)

type Util interface {
	Decipher(ctx context.Context, encryptedText string) (string, error)
}

func NewCryptoUtil(cryptoServiceUrl string) Util {
	return cryptoUtil{
		httpRequestBuilder: request.NewHttpRequestBuilder(util.GetHttpClientWithTracing()),
		cryptoServiceUrl:   cryptoServiceUrl,
	}
}

type cryptoUtil struct {
	httpRequestBuilder request.HttpRequestBuilder
	cryptoServiceUrl   string
}

func (utils cryptoUtil) Decipher(ctx context.Context, encryptedText string) (string, error) {
	if encryptedText == "" {
		return "", errors.New("text is empty")
	}

	response := &model.CryptoResponse{}

	url := utils.cryptoServiceUrl + constants.TEXT_DECRYPT_ROUTE
	err := utils.httpRequestBuilder.
		NewRequest().
		WithContext(ctx).
		WithJSONBody(model.CryptoRequest{EncryptedText: encryptedText}).
		ResponseAs(response).
		Post(url)
	if err != nil {
		return "", err
	}

	return response.DecryptedText, nil
}
