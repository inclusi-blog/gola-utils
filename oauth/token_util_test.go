package oauth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/constants"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/golang_error"
	"github.com/inclusi-blog/gola-utils/model"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type OauthUtilsTestSuite struct {
	suite.Suite
	recorder   *httptest.ResponseRecorder
	context    *gin.Context
	oauthUtils Utils
}

func TestOauthUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(OauthUtilsTestSuite))
}

func (suite *OauthUtilsTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "dummyUrl", nil)
}

func (suite OauthUtilsTestSuite) TestShouldReturnErrorWhenEncIdTokenHeaderIsNotPresentInGinContext() {
	suite.oauthUtils = NewOauthUtils("random-url")

	_, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)
	expectedError := errors.New("no enc_id_token present")
	suite.Equal(expectedError, err)
}

func (suite OauthUtilsTestSuite) TestShouldReturnErrorWhenEncIdTokenIsNotPresentInGoContext() {
	suite.oauthUtils = NewOauthUtils("random-url")

	_, err := suite.oauthUtils.DecodeEncryptedIdToken(context.Background())
	expectedError := errors.New("no enc_id_token present")
	suite.Equal(expectedError, err)
}

func (suite OauthUtilsTestSuite) TestShouldReturnErrorWhenAccessTokenIsNotPresentInGinContext() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.j5dPteSgV4TCsgTm5pL25auztQRJ1Eitofll7K84OKTcGx3RTHFN9dtD7_J8JSzXejCe3_LX3tHWQFykBUae8ttmMS9QGDM7nP0BrNserViN-_yM-CxLqz8E6BycE8PILuqLCqgGtyW2_QeTC26xESNpcCaVY2YHNyclTcgBWjGc3fxGGf5KI42yne8zRcnNtdQa5FTyQ-UG95T65ohROdKLRRGT12E9wGdEbcIj6ZJKeNt2fgL6eUFzPr8ReJwczWhsRzEk8E5yX-tJNcmdFKnicUcTIb_5PqfUIeaYNlCtS__NQP3wO8_kFIGiWv1fQslM77EdMe0CGy5g0DCtqw.TDXEYt_xckSm95GW.Penh8oBWbokLd80JqppLZkRuCLMh_UHg1WUi5L0TJz7y0x8ynaYnNgyD2tRSZJ5eOIjNam-PnzvyohabDFhQTN7oMcb1Y1kpToU588Ycxvd2a5redw_J8tPRdsNsAPDE2VU5bBaORsvHwUNiwbv6AxnL2s2E5EKGn9alO3Bmxu2VgVNtxjKtOh3Z0rfw5X6Lq9c-7yxyxT2hePXASSeXcypJGKHZd8AcIIDQZAFvgrpt6X7AfC4uby6TlR0d79FLdhMaW6p03tatG6T4HJRfbTASf-Xy1-SgslvXvh8ovtNKdXuGb7LT0M0i4fEQ4wIKS93gg3Fmv1eP9uN84MBGSRyExK5IlSD0VLpA4hlFnp98Dts3T3adGwFNaF6eurt4Bw6-TETXBi6crzD25lKUoQ34IcIUCgmrr1u60m-qsus7qfzP54h9OQ.Rg1WLjfSv_uwgdq4Ci0R4A"
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)
	suite.oauthUtils = NewOauthUtils("random-url")

	idToken, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)
	suite.Equal(model.IdToken{}, idToken)
	suite.Equal(errors.New("no access token present"), err)
}

func (suite OauthUtilsTestSuite) TestShouldReturnErrorWhenAccessTokenIsNotPresentInGoContext() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.j5dPteSgV4TCsgTm5pL25auztQRJ1Eitofll7K84OKTcGx3RTHFN9dtD7_J8JSzXejCe3_LX3tHWQFykBUae8ttmMS9QGDM7nP0BrNserViN-_yM-CxLqz8E6BycE8PILuqLCqgGtyW2_QeTC26xESNpcCaVY2YHNyclTcgBWjGc3fxGGf5KI42yne8zRcnNtdQa5FTyQ-UG95T65ohROdKLRRGT12E9wGdEbcIj6ZJKeNt2fgL6eUFzPr8ReJwczWhsRzEk8E5yX-tJNcmdFKnicUcTIb_5PqfUIeaYNlCtS__NQP3wO8_kFIGiWv1fQslM77EdMe0CGy5g0DCtqw.TDXEYt_xckSm95GW.Penh8oBWbokLd80JqppLZkRuCLMh_UHg1WUi5L0TJz7y0x8ynaYnNgyD2tRSZJ5eOIjNam-PnzvyohabDFhQTN7oMcb1Y1kpToU588Ycxvd2a5redw_J8tPRdsNsAPDE2VU5bBaORsvHwUNiwbv6AxnL2s2E5EKGn9alO3Bmxu2VgVNtxjKtOh3Z0rfw5X6Lq9c-7yxyxT2hePXASSeXcypJGKHZd8AcIIDQZAFvgrpt6X7AfC4uby6TlR0d79FLdhMaW6p03tatG6T4HJRfbTASf-Xy1-SgslvXvh8ovtNKdXuGb7LT0M0i4fEQ4wIKS93gg3Fmv1eP9uN84MBGSRyExK5IlSD0VLpA4hlFnp98Dts3T3adGwFNaF6eurt4Bw6-TETXBi6crzD25lKUoQ34IcIUCgmrr1u60m-qsus7qfzP54h9OQ.Rg1WLjfSv_uwgdq4Ci0R4A"
	suite.oauthUtils = NewOauthUtils("random-url")
	idToken, err := suite.oauthUtils.DecodeEncryptedIdToken(context.WithValue(context.Background(), constants.CONTEXT_ENC_ID_TOKEN, jweToken))
	suite.Equal(model.IdToken{}, idToken)
	suite.Equal(errors.New("no access token present"), err)
}

//TODO Change the key as enc_id_token to make it consistent. The change needs to done here and in crypto service.
func (suite OauthUtilsTestSuite) TestShouldReturnIdTokenForGivenValidJWETokenInGinContext() {
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVzZXJuYW1lIjoiZHVtbXktdXNlciIsImVtYWlsIjoiZHVtbXlAZ21haWwuY29tIiwic3ViamVjdCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsImF0X2hhc2giOiJ6NHJWZVY5SURXNG1RQmdrODBFUFdnIn0.IBLXYhF5TUUqyDsk5y2IzndElxhygZwRkuBYKmanBpU"
	suite.context.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer Glp4z8uis5TtrPXGDl2FQKYqNeaHgrhFlSLlLKvEd7U.ZsQVxn71lRB2QGHEJ94bS2aAKxUrSMj1CRMM6PtC4MY")
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jwtToken)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "id_token",
			Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVzZXJuYW1lIjoiZHVtbXktdXNlciIsImVtYWlsIjoiZHVtbXlAZ21haWwuY29tIiwic3ViamVjdCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsImF0X2hhc2giOiJ6NHJWZVY5SURXNG1RQmdrODBFUFdnIn0.IBLXYhF5TUUqyDsk5y2IzndElxhygZwRkuBYKmanBpU",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	idToken, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)
	expectedToken := model.IdToken{
		UserId:          "620943a2-ba51-4231-85b0-cc119c10840a",
		Username:        "dummy-user",
		Email:           "dummy@gmail.com",
		Subject:         "620943a2-ba51-4231-85b0-cc119c10840a",
		AccessTokenHash: "z4rVeV9IDW4mQBgk80EPWg",
	}
	suite.Equal(expectedToken, idToken)
	suite.Nil(err)
}

func (suite OauthUtilsTestSuite) TestShouldReturnIdTokenForGivenValidJWETokenInGoContext() {
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVzZXJuYW1lIjoiZHVtbXktdXNlciIsImVtYWlsIjoiZHVtbXlAZ21haWwuY29tIiwic3ViamVjdCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsImF0X2hhc2giOiJ6NHJWZVY5SURXNG1RQmdrODBFUFdnIn0.IBLXYhF5TUUqyDsk5y2IzndElxhygZwRkuBYKmanBpU"

	ctx := context.WithValue(context.Background(), constants.CONTEXT_ENC_ID_TOKEN, jwtToken)
	ctx = context.WithValue(ctx, constants.CONTEXT_ACCESS_TOKEN, "Glp4z8uis5TtrPXGDl2FQKYqNeaHgrhFlSLlLKvEd7U.ZsQVxn71lRB2QGHEJ94bS2aAKxUrSMj1CRMM6PtC4MY")

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "id_token",
			Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVzZXJuYW1lIjoiZHVtbXktdXNlciIsImVtYWlsIjoiZHVtbXlAZ21haWwuY29tIiwic3ViamVjdCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsImF0X2hhc2giOiJ6NHJWZVY5SURXNG1RQmdrODBFUFdnIn0.IBLXYhF5TUUqyDsk5y2IzndElxhygZwRkuBYKmanBpU",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	idToken, err := suite.oauthUtils.DecodeEncryptedIdToken(ctx)
	expectedToken := model.IdToken{
		UserId:          "620943a2-ba51-4231-85b0-cc119c10840a",
		Username:        "dummy-user",
		Email:           "dummy@gmail.com",
		Subject:         "620943a2-ba51-4231-85b0-cc119c10840a",
		AccessTokenHash: "z4rVeV9IDW4mQBgk80EPWg",
	}
	suite.Equal(expectedToken, idToken)
	suite.Nil(err)
}

//TODO Change the key as enc_id_token to make it consistent. The change needs to done here and in crypto service.
func (suite OauthUtilsTestSuite) TestShouldReturnErrorWhenProvidedIDTokenIsNotAssociatedWithAccessToken() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.j5dPteSgV4TCsgTm5pL25auztQRJ1Eitofll7K84OKTcGx3RTHFN9dtD7_J8JSzXejCe3_LX3tHWQFykBUae8ttmMS9QGDM7nP0BrNserViN-_yM-CxLqz8E6BycE8PILuqLCqgGtyW2_QeTC26xESNpcCaVY2YHNyclTcgBWjGc3fxGGf5KI42yne8zRcnNtdQa5FTyQ-UG95T65ohROdKLRRGT12E9wGdEbcIj6ZJKeNt2fgL6eUFzPr8ReJwczWhsRzEk8E5yX-tJNcmdFKnicUcTIb_5PqfUIeaYNlCtS__NQP3wO8_kFIGiWv1fQslM77EdMe0CGy5g0DCtqw.TDXEYt_xckSm95GW.Penh8oBWbokLd80JqppLZkRuCLMh_UHg1WUi5L0TJz7y0x8ynaYnNgyD2tRSZJ5eOIjNam-PnzvyohabDFhQTN7oMcb1Y1kpToU588Ycxvd2a5redw_J8tPRdsNsAPDE2VU5bBaORsvHwUNiwbv6AxnL2s2E5EKGn9alO3Bmxu2VgVNtxjKtOh3Z0rfw5X6Lq9c-7yxyxT2hePXASSeXcypJGKHZd8AcIIDQZAFvgrpt6X7AfC4uby6TlR0d79FLdhMaW6p03tatG6T4HJRfbTASf-Xy1-SgslvXvh8ovtNKdXuGb7LT0M0i4fEQ4wIKS93gg3Fmv1eP9uN84MBGSRyExK5IlSD0VLpA4hlFnp98Dts3T3adGwFNaF6eurt4Bw6-TETXBi6crzD25lKUoQ34IcIUCgmrr1u60m-qsus7qfzP54h9OQ.Rg1WLjfSv_uwgdq4Ci0R4A"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "id_token",
			Value: "eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzo1MzA4NWFmOS1lNGRmLTRlZWQtOThjZS1kMmNhYWE4MDQ0MzQiLCJ0eXAiOiJKV1QifQ.eyJhY3IiOiIxIiwiYXRfaGFzaCI6Ino0clZlVjlJRFc0bVFCZ2s4MEVQV2ciLCJhdWQiOlsiYTI0ZjA3MTYtYWE5MC00ZTU5LWJhNjAtZTliOTM5MjM3OWUzIl0sImF1dGhfdGltZSI6MTU4NDE4MTI3OSwiZW1haWwiOiJhbWl0LnlhZGF2QHF1YWxpdHlraW9zay5jb20iLCJleHAiOjE1ODQxODQ4OTAsImZpcnN0TmFtZSI6IlJhaHVsIiwiaWF0IjoxNTg0MTgxMjkwLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjQ0NDQvIiwianRpIjoiNGFhYTY4NzYtOTVkOC00NzExLWFlOTAtOTAxODViN2JmZTE2IiwibGFzdE5hbWUiOiJTaGVsa2UiLCJsYXN0U3VjY2Vzc2Z1bExvZ2luIjoiMjAyMC0wMy0xM1QxNjozMDowNC45NTY1MjJaIiwibG9naW5Nb2RlIjoiIiwibWFrZXJJZCI6IiIsIm1pZGRsZU5hbWUiOiJSYW5nbmF0aCIsIm1vYmlsZU51bWJlciI6IjkxOTY5NTQ4OTcwNCIsIm5vbmNlIjoiIiwicmF0IjoxNTg0MTgxMjUxLCJyZWxhdGlvbnMiOm51bGwsInNpZCI6Ijk1MmRjMTlkLWM1YTctNDk3OS1iMjhjLTQ5NWY2MTMzNWJiYSIsInN1YiI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVjaWMiOiI1MDAwMTU2MDUzIiwidXNlcklkIjoiNjIwOTQzYTItYmE1MS00MjMxLTg1YjAtY2MxMTljMTA4NDBhIn0.sNmrcQf4eHAnkrFDBAU_ZPmDT4CTMuMrsUHzZPmSbSIf9rvCO3EL3IBXg_cIYLzxAnRJ72Rz5yiAttEFdIlVu85IMmO6_OIGjB_M1WpvMGnqRVD9KBoiFPQ4NM9WiIK3V3eB17d-q2eHkpg4F4wZHT12qRlpRjYy9QVePRCBG8ictCD7zZI6RF3zw6nDf65pRPMx_zwQbHcGiAvtuhxPaa8kQEhIM261Co0zSVtHL2PYloXzLjifCGv80J4ifjDJxrvXbkQl02CoiUrGkViaMwvBkrQ-6A0rSeDqS1TvQDG7QOXQR6UJ6iOADvD5m1RuLQAwMs3ZxBOl_fnJELovlA6COA9SjkhAJDJ4Qg_6H4t3cENdhB-auogK8YBh81Gcn_HtRYpbDTKAVCy5VDeLH5W8sbimjbeiZvwMDH9QPS2hdwA_RRvYPQqH5PQvIOJkOmgEIQsTrK0GVNEY-mKG7I-bkXfE87FEKAG9h1yxFp8b8KYGn0F2EZ4pIQxI9-kA7ChNTwTd8vQt0vMmbP4Ay6GoljUBNprvmhrmOvY_m_4u1o9wZpc2RmDnh3S2x0-fmtNz07GEzDWctIJWSiixX6ghQjMkdRJytnBlz5_IqTlHH45qEOUNjJcfxKAUhSNtZceUuRn0FkVNu_8OIB80Uwk41eSQFgKWBLOqWwFL5F8",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)
	suite.context.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer Glp4z8uis5TtrPXGDl2FQKYqNeaHgrhFlSLlLKvEd7U.ZsQVxn71lRB2QGHEJ94bS2aAKxUrSMj1CRMM6PtC4MM")
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)

	expectedErr := golaerror.Error{ErrorCode: "ERR_UNAUTHORIZED", ErrorMessage: "ID token not associated with Access token"}
	idToken, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)
	suite.Equal(model.IdToken{}, idToken)
	suite.Equal(expectedErr, err)
}

func (suite OauthUtilsTestSuite) TestShouldReturnInternalServerErrorOnServiceFailure() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.F-yXHFhp3x-E1wJbsnt3LUbsL-wchaxKLKwMt6giCearkdyAsu3wAM8g6T2_QpTZtbCgXybBYEiog5pQY23NUVJwzPqh7KYakEHTm0UNjA-poyEEBJZt0EaFtjA2kPIvWeHg4pUqWxnIdLx2g_F4e8He9FXaCtESSiaSPJq7eir8ZaVxwqhflii3JhDS1pHRHjP7udxI-147SHUjISVYiWZgVu8MILfnh0IFsFTeGdrceo8ttCHNTKv7_ZUHOCchE94EMtPN3l-lkbq7l6hi7N6vIs4tLiyy_ri8beNVcOHoqlpmiNRPpH8-V0FAZioEQuB-ubrEQI1XB0bZIp_prg.4va3Db5iTeJ3FsS9.As6xdHXPw85XwNeRes3JYRZDFLq0KY93BYu3brSWRTJBH-Rks2GHL8e_wxxX-KZsJH7luMxKeZyhbFlN-iFHLuCQnWujjTvRii4CAelTCiYWYWQTeydaVpWGu7ANjKrJzk2ZEQ3jfM_bgbFulWeACa3mrcGED1G1Cuk3XVHcF-tf2MNiuc7TUyVFT-QmmNbNqSDdLz0v4vdxRY4axYpG1THVFJJ4Ut1Bq4vDb5H86uldthtrFSFNpzOL8DgCwvmnINLQJ4ml7yyk0EfbzMtJbr5G6x4d2XfnXyQw5kgGKNTWaf5l.5xHoBrEHj14ew5pG6aERfQ"
	suite.context.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer Glp4z8uis5TtrPXGDl2FQKYqNeaHgrhFlSLlLKvEd7U.ZsQVxn71lRB2QGHEJ94bS2aAKxUrSMj1CRMM6PtC4MY")
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(400)
		_, _ = res.Write([]byte(""))
	}))
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	defer func() { testServer.Close() }()
	_, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)

	suite.Equal(golang_error.InternalServerError{}, err)

}

func (suite OauthUtilsTestSuite) TestShouldReturnErrorForIncorrectIdToken() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.F-yXHFhp3x-E1wJbsnt3LUbsL-wchaxKLKwMt6giCearkdyAsu3wAM8g6T2_QpTZtbCgXybBYEiog5pQY23NUVJwzPqh7KYakEHTm0UNjA-poyEEBJZt0EaFtjA2kPIvWeHg4pUqWxnIdLx2g_F4e8He9FXaCtESSiaSPJq7eir8ZaVxwqhflii3JhDS1pHRHjP7udxI-147SHUjISVYiWZgVu8MILfnh0IFsFTeGdrceo8ttCHNTKv7_ZUHOCchE94EMtPN3l-lkbq7l6hi7N6vIs4tLiyy_ri8beNVcOHoqlpmiNRPpH8-V0FAZioEQuB-ubrEQI1XB0bZIp_prg.4va3Db5iTeJ3FsS9.As6xdHXPw85XwNeRes3JYRZDFLq0KY93BYu3brSWRTJBH-Rks2GHL8e_wxxX-KZsJH7luMxKeZyhbFlN-iFHLuCQnWujjTvRii4CAelTCiYWYWQTeydaVpWGu7ANjKrJzk2ZEQ3jfM_bgbFulWeACa3mrcGED1G1Cuk3XVHcF-tf2MNiuc7TUyVFT-QmmNbNqSDdLz0v4vdxRY4axYpG1THVFJJ4Ut1Bq4vDb5H86uldthtrFSFNpzOL8DgCwvmnINLQJ4ml7yyk0EfbzMtJbr5G6x4d2XfnXyQw5kgGKNTWaf5l.5xHoBrEHj14ew5pG6aERfQ"
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  constants.CONTEXT_ENC_ID_TOKEN,
			Value: "a.b.c",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)
	suite.NotNil(err)
}

func (suite OauthUtilsTestSuite) TestShouldReturnErrorWhenIDTokenIsNotThere() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.F-yXHFhp3x-E1wJbsnt3LUbsL-wchaxKLKwMt6giCearkdyAsu3wAM8g6T2_QpTZtbCgXybBYEiog5pQY23NUVJwzPqh7KYakEHTm0UNjA-poyEEBJZt0EaFtjA2kPIvWeHg4pUqWxnIdLx2g_F4e8He9FXaCtESSiaSPJq7eir8ZaVxwqhflii3JhDS1pHRHjP7udxI-147SHUjISVYiWZgVu8MILfnh0IFsFTeGdrceo8ttCHNTKv7_ZUHOCchE94EMtPN3l-lkbq7l6hi7N6vIs4tLiyy_ri8beNVcOHoqlpmiNRPpH8-V0FAZioEQuB-ubrEQI1XB0bZIp_prg.4va3Db5iTeJ3FsS9.As6xdHXPw85XwNeRes3JYRZDFLq0KY93BYu3brSWRTJBH-Rks2GHL8e_wxxX-KZsJH7luMxKeZyhbFlN-iFHLuCQnWujjTvRii4CAelTCiYWYWQTeydaVpWGu7ANjKrJzk2ZEQ3jfM_bgbFulWeACa3mrcGED1G1Cuk3XVHcF-tf2MNiuc7TUyVFT-QmmNbNqSDdLz0v4vdxRY4axYpG1THVFJJ4Ut1Bq4vDb5H86uldthtrFSFNpzOL8DgCwvmnINLQJ4ml7yyk0EfbzMtJbr5G6x4d2XfnXyQw5kgGKNTWaf5l.5xHoBrEHj14ew5pG6aERfQ"
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)

	suite.NotNil(err)
}

func (suite OauthUtilsTestSuite) TestDecodeEncryptedIdToken_ShouldReturnErrorWhenCryptoServiceDoesNotReturnIDTokenCookie() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.j5dPteSgV4TCsgTm5pL25auztQRJ1Eitofll7K84OKTcGx3RTHFN9dtD7_J8JSzXejCe3_LX3tHWQFykBUae8ttmMS9QGDM7nP0BrNserViN-_yM-CxLqz8E6BycE8PILuqLCqgGtyW2_QeTC26xESNpcCaVY2YHNyclTcgBWjGc3fxGGf5KI42yne8zRcnNtdQa5FTyQ-UG95T65ohROdKLRRGT12E9wGdEbcIj6ZJKeNt2fgL6eUFzPr8ReJwczWhsRzEk8E5yX-tJNcmdFKnicUcTIb_5PqfUIeaYNlCtS__NQP3wO8_kFIGiWv1fQslM77EdMe0CGy5g0DCtqw.TDXEYt_xckSm95GW.Penh8oBWbokLd80JqppLZkRuCLMh_UHg1WUi5L0TJz7y0x8ynaYnNgyD2tRSZJ5eOIjNam-PnzvyohabDFhQTN7oMcb1Y1kpToU588Ycxvd2a5redw_J8tPRdsNsAPDE2VU5bBaORsvHwUNiwbv6AxnL2s2E5EKGn9alO3Bmxu2VgVNtxjKtOh3Z0rfw5X6Lq9c-7yxyxT2hePXASSeXcypJGKHZd8AcIIDQZAFvgrpt6X7AfC4uby6TlR0d79FLdhMaW6p03tatG6T4HJRfbTASf-Xy1-SgslvXvh8ovtNKdXuGb7LT0M0i4fEQ4wIKS93gg3Fmv1eP9uN84MBGSRyExK5IlSD0VLpA4hlFnp98Dts3T3adGwFNaF6eurt4Bw6-TETXBi6crzD25lKUoQ34IcIUCgmrr1u60m-qsus7qfzP54h9OQ.Rg1WLjfSv_uwgdq4Ci0R4A"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "id_TOken",
			Value: "eyJhbGciOiJSUzI1NiIsImtpZCI6InB1YmxpYzo1MzA4NWFmOS1lNGRmLTRlZWQtOThjZS1kMmNhYWE4MDQ0MzQiLCJ0eXAiOiJKV1QifQ.eyJhY3IiOiIxIiwiYXRfaGFzaCI6Ino0clZlVjlJRFc0bVFCZ2s4MEVQV2ciLCJhdWQiOlsiYTI0ZjA3MTYtYWE5MC00ZTU5LWJhNjAtZTliOTM5MjM3OWUzIl0sImF1dGhfdGltZSI6MTU4NDE4MTI3OSwiZW1haWwiOiJhbWl0LnlhZGF2QHF1YWxpdHlraW9zay5jb20iLCJleHAiOjE1ODQxODQ4OTAsImZpcnN0TmFtZSI6IlJhaHVsIiwiaWF0IjoxNTg0MTgxMjkwLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjQ0NDQvIiwianRpIjoiNGFhYTY4NzYtOTVkOC00NzExLWFlOTAtOTAxODViN2JmZTE2IiwibGFzdE5hbWUiOiJTaGVsa2UiLCJsYXN0U3VjY2Vzc2Z1bExvZ2luIjoiMjAyMC0wMy0xM1QxNjozMDowNC45NTY1MjJaIiwibG9naW5Nb2RlIjoiIiwibWFrZXJJZCI6IiIsIm1pZGRsZU5hbWUiOiJSYW5nbmF0aCIsIm1vYmlsZU51bWJlciI6IjkxOTY5NTQ4OTcwNCIsIm5vbmNlIjoiIiwicmF0IjoxNTg0MTgxMjUxLCJyZWxhdGlvbnMiOm51bGwsInNpZCI6Ijk1MmRjMTlkLWM1YTctNDk3OS1iMjhjLTQ5NWY2MTMzNWJiYSIsInN1YiI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVjaWMiOiI1MDAwMTU2MDUzIiwidXNlcklkIjoiNjIwOTQzYTItYmE1MS00MjMxLTg1YjAtY2MxMTljMTA4NDBhIn0.sNmrcQf4eHAnkrFDBAU_ZPmDT4CTMuMrsUHzZPmSbSIf9rvCO3EL3IBXg_cIYLzxAnRJ72Rz5yiAttEFdIlVu85IMmO6_OIGjB_M1WpvMGnqRVD9KBoiFPQ4NM9WiIK3V3eB17d-q2eHkpg4F4wZHT12qRlpRjYy9QVePRCBG8ictCD7zZI6RF3zw6nDf65pRPMx_zwQbHcGiAvtuhxPaa8kQEhIM261Co0zSVtHL2PYloXzLjifCGv80J4ifjDJxrvXbkQl02CoiUrGkViaMwvBkrQ-6A0rSeDqS1TvQDG7QOXQR6UJ6iOADvD5m1RuLQAwMs3ZxBOl_fnJELovlA6COA9SjkhAJDJ4Qg_6H4t3cENdhB-auogK8YBh81Gcn_HtRYpbDTKAVCy5VDeLH5W8sbimjbeiZvwMDH9QPS2hdwA_RRvYPQqH5PQvIOJkOmgEIQsTrK0GVNEY-mKG7I-bkXfE87FEKAG9h1yxFp8b8KYGn0F2EZ4pIQxI9-kA7ChNTwTd8vQt0vMmbP4Ay6GoljUBNprvmhrmOvY_m_4u1o9wZpc2RmDnh3S2x0-fmtNz07GEzDWctIJWSiixX6ghQjMkdRJytnBlz5_IqTlHH45qEOUNjJcfxKAUhSNtZceUuRn0FkVNu_8OIB80Uwk41eSQFgKWBLOqWwFL5F8",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.context.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer Glp4z8uis5TtrPXGDl2FQKYqNeaHgrhFlSLlLKvEd7U.ZsQVxn71lRB2QGHEJ94bS2aAKxUrSMj1CRMM6PtC4MY")
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)

	suite.Equal(errors.New("cookie not found"), err)
}

//TODO Change the key as enc_id_token to make it consistent. The change needs to done here and in crypto service.
func (suite OauthUtilsTestSuite) TestDecodeEncryptedIdToken_ShouldReturnErrorIDTokenIsNotInTheRightFormat() {
	jweToken := "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.j5dPteSgV4TCsgTm5pL25auztQRJ1Eitofll7K84OKTcGx3RTHFN9dtD7_J8JSzXejCe3_LX3tHWQFykBUae8ttmMS9QGDM7nP0BrNserViN-_yM-CxLqz8E6BycE8PILuqLCqgGtyW2_QeTC26xESNpcCaVY2YHNyclTcgBWjGc3fxGGf5KI42yne8zRcnNtdQa5FTyQ-UG95T65ohROdKLRRGT12E9wGdEbcIj6ZJKeNt2fgL6eUFzPr8ReJwczWhsRzEk8E5yX-tJNcmdFKnicUcTIb_5PqfUIeaYNlCtS__NQP3wO8_kFIGiWv1fQslM77EdMe0CGy5g0DCtqw.TDXEYt_xckSm95GW.Penh8oBWbokLd80JqppLZkRuCLMh_UHg1WUi5L0TJz7y0x8ynaYnNgyD2tRSZJ5eOIjNam-PnzvyohabDFhQTN7oMcb1Y1kpToU588Ycxvd2a5redw_J8tPRdsNsAPDE2VU5bBaORsvHwUNiwbv6AxnL2s2E5EKGn9alO3Bmxu2VgVNtxjKtOh3Z0rfw5X6Lq9c-7yxyxT2hePXASSeXcypJGKHZd8AcIIDQZAFvgrpt6X7AfC4uby6TlR0d79FLdhMaW6p03tatG6T4HJRfbTASf-Xy1-SgslvXvh8ovtNKdXuGb7LT0M0i4fEQ4wIKS93gg3Fmv1eP9uN84MBGSRyExK5IlSD0VLpA4hlFnp98Dts3T3adGwFNaF6eurt4Bw6-TETXBi6crzD25lKUoQ34IcIUCgmrr1u60m-qsus7qfzP54h9OQ.Rg1WLjfSv_uwgdq4Ci0R4A"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "id_token",
			Value: "eyJhY3IiOiIxIiwiYXRfaGFzaCI6Ino0clZlVjlJRFc0bVFCZ2s4MEVQV2ciLCJhdWQiOlsiYTI0ZjA3MTYtYWE5MC00ZTU5LWJhNjAtZTliOTM5MjM3OWUzIl0sImF1dGhfdGltZSI6MTU4NDE4MTI3OSwiZW1haWwiOiJhbWl0LnlhZGF2QHF1YWxpdHlraW9zay5jb20iLCJleHAiOjE1ODQxODQ4OTAsImZpcnN0TmFtZSI6IlJhaHVsIiwiaWF0IjoxNTg0MTgxMjkwLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjQ0NDQvIiwianRpIjoiNGFhYTY4NzYtOTVkOC00NzExLWFlOTAtOTAxODViN2JmZTE2IiwibGFzdE5hbWUiOiJTaGVsa2UiLCJsYXN0U3VjY2Vzc2Z1bExvZ2luIjoiMjAyMC0wMy0xM1QxNjozMDowNC45NTY1MjJaIiwibG9naW5Nb2RlIjoiIiwibWFrZXJJZCI6IiIsIm1pZGRsZU5hbWUiOiJSYW5nbmF0aCIsIm1vYmlsZU51bWJlciI6IjkxOTY5NTQ4OTcwNCIsIm5vbmNlIjoiIiwicmF0IjoxNTg0MTgxMjUxLCJyZWxhdGlvbnMiOm51bGwsInNpZCI6Ijk1MmRjMTlkLWM1YTctNDk3OS1iMjhjLTQ5NWY2MTMzNWJiYSIsInN1YiI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVjaWMiOiI1MDAwMTU2MDUzIiwidXNlcklkIjoiNjIwOTQzYTItYmE1MS00MjMxLTg1YjAtY2MxMTljMTA4NDBhIn0.sNmrcQf4eHAnkrFDBAU_ZPmDT4CTMuMrsUHzZPmSbSIf9rvCO3EL3IBXg_cIYLzxAnRJ72Rz5yiAttEFdIlVu85IMmO6_OIGjB_M1WpvMGnqRVD9KBoiFPQ4NM9WiIK3V3eB17d-q2eHkpg4F4wZHT12qRlpRjYy9QVePRCBG8ictCD7zZI6RF3zw6nDf65pRPMx_zwQbHcGiAvtuhxPaa8kQEhIM261Co0zSVtHL2PYloXzLjifCGv80J4ifjDJxrvXbkQl02CoiUrGkViaMwvBkrQ-6A0rSeDqS1TvQDG7QOXQR6UJ6iOADvD5m1RuLQAwMs3ZxBOl_fnJELovlA6COA9SjkhAJDJ4Qg_6H4t3cENdhB-auogK8YBh81Gcn_HtRYpbDTKAVCy5VDeLH5W8sbimjbeiZvwMDH9QPS2hdwA_RRvYPQqH5PQvIOJkOmgEIQsTrK0GVNEY-mKG7I-bkXfE87FEKAG9h1yxFp8b8KYGn0F2EZ4pIQxI9-kA7ChNTwTd8vQt0vMmbP4Ay6GoljUBNprvmhrmOvY_m_4u1o9wZpc2RmDnh3S2x0-fmtNz07GEzDWctIJWSiixX6ghQjMkdRJytnBlz5_IqTlHH45qEOUNjJcfxKAUhSNtZceUuRn0FkVNu_8OIB80Uwk41eSQFgKWBLOqWwFL5F8",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.context.Request.Header.Set(constants.AUTHORIZATION_HEADER_KEY, "Bearer Glp4z8uis5TtrPXGDl2FQKYqNeaHgrhFlSLlLKvEd7U.ZsQVxn71lRB2QGHEJ94bS2aAKxUrSMj1CRMM6PtC4MY")
	suite.context.Request.Header.Set(constants.ENC_ID_TOKEN_HEADER_KEY, jweToken)
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.DecodeEncryptedIdToken(suite.context)

	suite.Equal(errors.New("invalid token format"), err)
}

func (suite OauthUtilsTestSuite) TestEncryptIdTokenShouldReturnValidJWEWhenIdTokenIsProvided() {
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImZpcnN0TmFtZSI6IiIsImxhc3ROYW1lIjoiIiwibGFzdFN1Y2Nlc3NmdWxMb2dpbiI6IjIwMTktMTAtMjRUMTU6MTc6MjIrMDU6MzAiLCJtaWRkbGVOYW1lIjoiIiwibW9iaWxlTnVtYmVyIjoiIiwic3ViIjoiNmQ0MWU3YWQtZGM5NC00ODhhLThjZDgtNzk1ZGMzNmM2MTRjIiwidWNpYyI6IiIsInVzZXJJZCI6IjZkNDFlN2FkLWRjOTQtNDg4YS04Y2Q4LTc5NWRjMzZjNjE0YyJ9.NL7qzA-TSmKEOhmEvudBZx3ozOFIfCCNmamw3vNDCe3NdUCu9ba1vG2N2Oj79uPeTtnZqHPQDwBA7aUGARxFa3xjybEdtdHNI-UWwjLxJ3IKdkwojQz1TCQznr7uCB9PLCrAneEKj4KvC1laeE9aL-FfzO4oRY1aU2zV7lVRsPkvFTCSVYkRFAo0Yejtz-E2nhnIdRjQ3IwfWkEDjujohLxa16kOsmxEWTxTeFknoa1VoIleh67Vk8Z80KrUbjfc0NW2DZ_wjOrzFXVxoOFrZDPYYG_CIR9_uJjjHAW2daQlfrLvb-rvb1FU1h1gGSgg_hPOHTwbJzM6flykiLkW3w"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "enc_id_token",
			Value: "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.F-yXHFhp3x-E1wJbsnt3LUbsL-wchaxKLKwMt6giCearkdyAsu3wAM8g6T2_QpTZtbCgXybBYEiog5pQY23NUVJwzPqh7KYakEHTm0UNjA-poyEEBJZt0EaFtjA2kPIvWeHg4pUqWxnIdLx2g_F4e8He9FXaCtESSiaSPJq7eir8ZaVxwqhflii3JhDS1pHRHjP7udxI-147SHUjISVYiWZgVu8MILfnh0IFsFTeGdrceo8ttCHNTKv7_ZUHOCchE94EMtPN3l-lkbq7l6hi7N6vIs4tLiyy_ri8beNVcOHoqlpmiNRPpH8-V0FAZioEQuB-ubrEQI1XB0bZIp_prg.4va3Db5iTeJ3FsS9.As6xdHXPw85XwNeRes3JYRZDFLq0KY93BYu3brSWRTJBH-Rks2GHL8e_wxxX-KZsJH7luMxKeZyhbFlN-iFHLuCQnWujjTvRii4CAelTCiYWYWQTeydaVpWGu7ANjKrJzk2ZEQ3jfM_bgbFulWeACa3mrcGED1G1Cuk3XVHcF-tf2MNiuc7TUyVFT-QmmNbNqSDdLz0v4vdxRY4axYpG1THVFJJ4Ut1Bq4vDb5H86uldthtrFSFNpzOL8DgCwvmnINLQJ4ml7yyk0EfbzMtJbr5G6x4d2XfnXyQw5kgGKNTWaf5l.5xHoBrEHj14ew5pG6aERfQ",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	jweToken, err := suite.oauthUtils.EncryptIdToken(suite.context, jwtToken)
	suite.Equal(5, len(strings.Split(jweToken, ".")))
	suite.Nil(err)
}

func (suite OauthUtilsTestSuite) TestEncryptIdTokenShouldReturnErrorWhenJWETokenNotFound() {
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImZpcnN0TmFtZSI6IiIsImxhc3ROYW1lIjoiIiwibGFzdFN1Y2Nlc3NmdWxMb2dpbiI6IjIwMTktMTAtMjRUMTU6MTc6MjIrMDU6MzAiLCJtaWRkbGVOYW1lIjoiIiwibW9iaWxlTnVtYmVyIjoiIiwic3ViIjoiNmQ0MWU3YWQtZGM5NC00ODhhLThjZDgtNzk1ZGMzNmM2MTRjIiwidWNpYyI6IiIsInVzZXJJZCI6IjZkNDFlN2FkLWRjOTQtNDg4YS04Y2Q4LTc5NWRjMzZjNjE0YyJ9.NL7qzA-TSmKEOhmEvudBZx3ozOFIfCCNmamw3vNDCe3NdUCu9ba1vG2N2Oj79uPeTtnZqHPQDwBA7aUGARxFa3xjybEdtdHNI-UWwjLxJ3IKdkwojQz1TCQznr7uCB9PLCrAneEKj4KvC1laeE9aL-FfzO4oRY1aU2zV7lVRsPkvFTCSVYkRFAo0Yejtz-E2nhnIdRjQ3IwfWkEDjujohLxa16kOsmxEWTxTeFknoa1VoIleh67Vk8Z80KrUbjfc0NW2DZ_wjOrzFXVxoOFrZDPYYG_CIR9_uJjjHAW2daQlfrLvb-rvb1FU1h1gGSgg_hPOHTwbJzM6flykiLkW3w"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.EncryptIdToken(suite.context, jwtToken)
	suite.NotNil(err)
}

func (suite OauthUtilsTestSuite) TestEncryptIdTokenShouldReturnErrorWhenJWTTokenIsNotProvidedForEncryption() {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "enc_id_token",
			Value: "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.F-yXHFhp3x-E1wJbsnt3LUbsL-wchaxKLKwMt6giCearkdyAsu3wAM8g6T2_QpTZtbCgXybBYEiog5pQY23NUVJwzPqh7KYakEHTm0UNjA-poyEEBJZt0EaFtjA2kPIvWeHg4pUqWxnIdLx2g_F4e8He9FXaCtESSiaSPJq7eir8ZaVxwqhflii3JhDS1pHRHjP7udxI-147SHUjISVYiWZgVu8MILfnh0IFsFTeGdrceo8ttCHNTKv7_ZUHOCchE94EMtPN3l-lkbq7l6hi7N6vIs4tLiyy_ri8beNVcOHoqlpmiNRPpH8-V0FAZioEQuB-ubrEQI1XB0bZIp_prg.4va3Db5iTeJ3FsS9.As6xdHXPw85XwNeRes3JYRZDFLq0KY93BYu3brSWRTJBH-Rks2GHL8e_wxxX-KZsJH7luMxKeZyhbFlN-iFHLuCQnWujjTvRii4CAelTCiYWYWQTeydaVpWGu7ANjKrJzk2ZEQ3jfM_bgbFulWeACa3mrcGED1G1Cuk3XVHcF-tf2MNiuc7TUyVFT-QmmNbNqSDdLz0v4vdxRY4axYpG1THVFJJ4Ut1Bq4vDb5H86uldthtrFSFNpzOL8DgCwvmnINLQJ4ml7yyk0EfbzMtJbr5G6x4d2XfnXyQw5kgGKNTWaf5l.5xHoBrEHj14ew5pG6aERfQ",
		})
		res.WriteHeader(200)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.EncryptIdToken(suite.context, "")
	suite.NotNil(err)
}

func (suite OauthUtilsTestSuite) TestEncryptIdTokenShouldReturnErrorWhenCryptoEndpointReturnsError() {
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImZpcnN0TmFtZSI6IiIsImxhc3ROYW1lIjoiIiwibGFzdFN1Y2Nlc3NmdWxMb2dpbiI6IjIwMTktMTAtMjRUMTU6MTc6MjIrMDU6MzAiLCJtaWRkbGVOYW1lIjoiIiwibW9iaWxlTnVtYmVyIjoiIiwic3ViIjoiNmQ0MWU3YWQtZGM5NC00ODhhLThjZDgtNzk1ZGMzNmM2MTRjIiwidWNpYyI6IiIsInVzZXJJZCI6IjZkNDFlN2FkLWRjOTQtNDg4YS04Y2Q4LTc5NWRjMzZjNjE0YyJ9.NL7qzA-TSmKEOhmEvudBZx3ozOFIfCCNmamw3vNDCe3NdUCu9ba1vG2N2Oj79uPeTtnZqHPQDwBA7aUGARxFa3xjybEdtdHNI-UWwjLxJ3IKdkwojQz1TCQznr7uCB9PLCrAneEKj4KvC1laeE9aL-FfzO4oRY1aU2zV7lVRsPkvFTCSVYkRFAo0Yejtz-E2nhnIdRjQ3IwfWkEDjujohLxa16kOsmxEWTxTeFknoa1VoIleh67Vk8Z80KrUbjfc0NW2DZ_wjOrzFXVxoOFrZDPYYG_CIR9_uJjjHAW2daQlfrLvb-rvb1FU1h1gGSgg_hPOHTwbJzM6flykiLkW3w"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.SetCookie(res, &http.Cookie{
			Name:  "enc_id_token",
			Value: "eyJhbGciOiJSU0EtT0FFUCIsImVuYyI6IkExMjhHQ00ifQ.F-yXHFhp3x-E1wJbsnt3LUbsL-wchaxKLKwMt6giCearkdyAsu3wAM8g6T2_QpTZtbCgXybBYEiog5pQY23NUVJwzPqh7KYakEHTm0UNjA-poyEEBJZt0EaFtjA2kPIvWeHg4pUqWxnIdLx2g_F4e8He9FXaCtESSiaSPJq7eir8ZaVxwqhflii3JhDS1pHRHjP7udxI-147SHUjISVYiWZgVu8MILfnh0IFsFTeGdrceo8ttCHNTKv7_ZUHOCchE94EMtPN3l-lkbq7l6hi7N6vIs4tLiyy_ri8beNVcOHoqlpmiNRPpH8-V0FAZioEQuB-ubrEQI1XB0bZIp_prg.4va3Db5iTeJ3FsS9.As6xdHXPw85XwNeRes3JYRZDFLq0KY93BYu3brSWRTJBH-Rks2GHL8e_wxxX-KZsJH7luMxKeZyhbFlN-iFHLuCQnWujjTvRii4CAelTCiYWYWQTeydaVpWGu7ANjKrJzk2ZEQ3jfM_bgbFulWeACa3mrcGED1G1Cuk3XVHcF-tf2MNiuc7TUyVFT-QmmNbNqSDdLz0v4vdxRY4axYpG1THVFJJ4Ut1Bq4vDb5H86uldthtrFSFNpzOL8DgCwvmnINLQJ4ml7yyk0EfbzMtJbr5G6x4d2XfnXyQw5kgGKNTWaf5l.5xHoBrEHj14ew5pG6aERfQ",
		})
		res.WriteHeader(500)
		_, _ = res.Write([]byte(""))
	}))
	defer func() { testServer.Close() }()
	suite.oauthUtils = NewOauthUtils(testServer.URL)

	_, err := suite.oauthUtils.EncryptIdToken(suite.context, jwtToken)
	suite.NotNil(err)
	suite.Equal(errors.New("something went wrong"), err)
}

func (suite OauthUtilsTestSuite) TestDecodeIdTokenFromJWTToken_ShouldReturnErrorWhenTokenInvalid() {
	jwtToken := "invalid.format"
	expectedError := errors.New("invalid token format")
	suite.oauthUtils = NewOauthUtils("http://crypto.servic-url")
	_, err := suite.oauthUtils.DecodeIdTokenFromJWT(jwtToken)
	suite.Equal(expectedError, err)
}

func (suite OauthUtilsTestSuite) TestDecodeIdTokenFromJWTToken_ShouldReturnSucessWhenTokenIsValid() {
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsInVzZXJuYW1lIjoiZHVtbXktdXNlciIsImVtYWlsIjoiZHVtbXlAZ21haWwuY29tIiwic3ViamVjdCI6IjYyMDk0M2EyLWJhNTEtNDIzMS04NWIwLWNjMTE5YzEwODQwYSIsImF0X2hhc2giOiJ6NHJWZVY5SURXNG1RQmdrODBFUFdnIn0.IBLXYhF5TUUqyDsk5y2IzndElxhygZwRkuBYKmanBpU"
	expectedToken := model.IdToken{
		UserId:          "620943a2-ba51-4231-85b0-cc119c10840a",
		Username:        "dummy-user",
		Email:           "dummy@gmail.com",
		Subject:         "620943a2-ba51-4231-85b0-cc119c10840a",
		AccessTokenHash: "z4rVeV9IDW4mQBgk80EPWg",
	}
	suite.oauthUtils = NewOauthUtils("http://crypto.servic-url")
	idToken, err := suite.oauthUtils.DecodeIdTokenFromJWT(jwtToken)
	suite.Nil(err)
	suite.Equal(expectedToken, idToken)
}

func (suite OauthUtilsTestSuite) TestDecodeIdTokenFromJWTToken_ShouldReturnErrorWhenUnmarshallingFails() {
	jwtToken := "ab.bc.cd"
	suite.oauthUtils = NewOauthUtils("http://crypto.servic-url")

	idToken, err := suite.oauthUtils.DecodeIdTokenFromJWT(jwtToken)

	suite.NotNil(err)
	suite.Empty(idToken)
}
