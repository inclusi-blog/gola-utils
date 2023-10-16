package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/inclusi-blog/gola-utils/logging"
	"go.opencensus.io/plugin/ochttp"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetHttpClientWithTracing() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 50 * time.Second,
		}).DialContext,
	}
	return &http.Client{Transport: &ochttp.Transport{Base: transport}}
}

func getBearerTokenFromHeader(context *gin.Context) (string, error) {
	bearerToken := context.Request.Header.Get(constants.AUTHORIZATION_HEADER_KEY)
	if bearerToken == "" {
		return "", errors.New("empty bearer token/not present")
	}
	if strings.HasPrefix(bearerToken, "Bearer") && len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", errors.New("invalid bearer token")
}

func getAccessTokenFromCookie(context *gin.Context) (string, error) {
	return context.Cookie(constants.COOKIE_ACCESS_TOKEN)
}

func GetAccessToken(context *gin.Context) (string, error) {
	logger := logging.GetLogger(context)
	if accessToken, bearerTokenErr := getBearerTokenFromHeader(context); bearerTokenErr == nil {
		logger.Info("using access token from header")
		return accessToken, nil
	} else if accessToken, cookieErr := getAccessTokenFromCookie(context); cookieErr == nil {
		logger.Warn(bearerTokenErr.Error())
		logger.Info("using access token from cookie")
		return accessToken, nil
	} else {
		logger.Warn(cookieErr.Error())
		return "", errors.New("invalid bearer/cookie header")
	}
}

func GetEncryptedIDToken(context *gin.Context) (string, error) {
	logger := logging.GetLogger(context)
	if encIDToken := context.Request.Header.Get(constants.ENC_ID_TOKEN_HEADER_KEY); encIDToken != "" {
		logger.Info("using enc id-token from header")
		return encIDToken, nil
	} else if encIDToken, err := context.Cookie(constants.COOKIE_ENC_ID_TOKEN); err == nil {
		logger.Warn("encrypted ID token in header is empty / not present")
		logger.Info("using enc id-token from cookie")
		return encIDToken, nil
	} else {
		logger.Warn(err.Error())
		return "", errors.New("invalid enc-id-token header/cookie")
	}
}

func FormBearerAuthorizationHeader(accessToken string) string {
	return "Bearer " + accessToken
}
