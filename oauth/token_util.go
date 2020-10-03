package oauth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/constants"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/golang_error"
	"github.com/gola-glitch/gola-utils/http/request"
	"github.com/gola-glitch/gola-utils/http/util"
	"github.com/gola-glitch/gola-utils/model"
	"net/http"
	"reflect"
	"strings"
)

type Utils interface {
	DecodeEncryptedIdToken(ctx context.Context) (model.IdToken, error)
	EncryptIdToken(ctx *gin.Context, token string) (string, error)
	DecodeIdTokenFromJWT(idToken string) (model.IdToken, error)
}

func NewOauthUtils(cryptoServiceUrl string) Utils {
	return oauthUtils{
		httpRequestBuilder: request.NewHttpRequestBuilder(util.GetHttpClientWithTracing()),
		cryptoServiceUrl:   cryptoServiceUrl,
	}
}

type oauthUtils struct {
	httpRequestBuilder request.HttpRequestBuilder
	cryptoServiceUrl   string
}

func (utils oauthUtils) DecodeIdTokenFromJWT(idToken string) (model.IdToken, error) {
	tokenParts := strings.Split(idToken, ".")

	if len(tokenParts) != 3 {
		return model.IdToken{}, errors.New("invalid token format")
	}
	payload := tokenParts[1]
	segment, decodeErr := jwt.DecodeSegment(payload)
	if decodeErr != nil {
		return model.IdToken{}, decodeErr
	}
	var token model.IdToken
	unmarshalError := json.Unmarshal(segment, &token)
	if unmarshalError != nil {
		return model.IdToken{}, unmarshalError
	}
	return token, nil
}

func (utils oauthUtils) DecodeEncryptedIdToken(ctx context.Context) (model.IdToken, error) {
	var jweToken string
	var accessToken string
	if ginContext, ok := ctx.(*gin.Context); ok {
		jweTokenCookie, jweTokenErr := util.GetEncryptedIDToken(ginContext)
		if jweTokenErr != nil {
			return model.IdToken{}, errors.New("no enc_id_token present")
		}
		jweToken = jweTokenCookie

		accessTokenCookie, accessTokenErr := util.GetAccessToken(ginContext)
		if accessTokenErr != nil {
			return model.IdToken{}, errors.New("no access token present")
		}
		accessToken = accessTokenCookie
	} else {
		jweTokenCookie := ctx.Value(constants.CONTEXT_ENC_ID_TOKEN)
		if jweTokenCookie == nil {
			return model.IdToken{}, errors.New("no enc_id_token present")
		}
		jweToken = jweTokenCookie.(string)

		accessTokenCookie := ctx.Value(constants.CONTEXT_ACCESS_TOKEN)
		if accessTokenCookie == nil {
			return model.IdToken{}, errors.New("no access token present")
		}
		accessToken = accessTokenCookie.(string)
	}
	url := utils.cryptoServiceUrl + constants.TOKEN_DECRYPT_ROUTE

	var responseCookies = &[]*http.Cookie{}
	responseBody := ""
	err := utils.
		httpRequestBuilder.
		NewRequest().
		WithContext(ctx).
		AddCookie(&http.Cookie{Name: "enc_id_token", Value: jweToken}).
		ResponseAs(&responseBody).
		ResponseCookiesAs(responseCookies).
		Get(url)
	if err != nil {
		return model.IdToken{}, golang_error.InternalServerError{}
	}

	extractedIdToken, err := utils.extractIdToken(*responseCookies)
	if err != nil {
		return model.IdToken{}, err
	}

	return utils.ensureIDTokenIsAssociatedWithAccessToken(accessToken, extractedIdToken)
}

func (utils oauthUtils) EncryptIdToken(ctx *gin.Context, jwtToken string) (string, error) {
	if jwtToken == "" {
		return "", errors.New("no id_token present")
	}

	url := utils.cryptoServiceUrl + constants.TOKEN_ENCRYPT_ROUTE
	var responseCookies = &[]*http.Cookie{}
	responseBody := ""
	//TODO Change the key as enc_id_token to make it consistent. The change needs to done here and in crypto service.
	err := utils.
		httpRequestBuilder.
		NewRequest().
		WithContext(ctx).
		AddCookie(&http.Cookie{Name: "id_token", Value: jwtToken}).
		ResponseAs(&responseBody).
		ResponseCookiesAs(responseCookies).
		Get(url)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	return utils.extractCookie(*responseCookies, "enc_id_token")
}

func (utils oauthUtils) extractCookie(cookies []*http.Cookie, name string) (string, error) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value, nil
		}
	}

	return "", errors.New("cookie not found")
}

func (utils oauthUtils) extractIdToken(cookies []*http.Cookie) (model.IdToken, error) {
	//TODO Change the key as enc_id_token to make it consistent. The change needs to done here and in crypto service.
	tokenString, err := utils.extractCookie(cookies, "id_token")
	if err != nil {
		return model.IdToken{}, err
	}

	token, err := utils.DecodeIdTokenFromJWT(tokenString)
	if err != nil {
		return model.IdToken{}, err
	}

	if reflect.DeepEqual(token, model.IdToken{}) {
		return model.IdToken{}, errors.New("id_token not found")
	}
	return token, nil
}

func (utils oauthUtils) ensureIDTokenIsAssociatedWithAccessToken(accessToken string, idToken model.IdToken) (model.IdToken, error) {

	buffer := bytes.NewBufferString(accessToken)
	hash := sha256.New()
	hash.Write(buffer.Bytes())
	hashBuf := bytes.NewBuffer(hash.Sum([]byte{}))
	length := hashBuf.Len()

	if idToken.AccessTokenHash != base64.RawURLEncoding.EncodeToString(hashBuf.Bytes()[:length/2]) {
		return model.IdToken{}, golaerror.Error{ErrorCode: "ERR_UNAUTHORIZED", ErrorMessage: "ID token not associated with Access token"}
	}
	return idToken, nil
}
