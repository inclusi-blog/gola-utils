package email

import (
	"context"
	"encoding/json"
	"github.com/inclusi-blog/gola-utils/alert/email/models"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

type EmailUtilTest struct {
	suite.Suite
	emailUtil  Util
	gatewayURL string
	goContext  context.Context
}

func (suite *EmailUtilTest) SetupTest() {
	suite.gatewayURL = "http://test.com/sendEmail"
	suite.emailUtil = NewEmailUtil(suite.gatewayURL)
	suite.goContext = context.WithValue(context.TODO(), "testKey", "testVal")
}

func TestEmailUtilSuite(t *testing.T) {
	suite.Run(t, new(EmailUtilTest))
}

func (suite *EmailUtilTest) TestShouldSendEmail_WhenValidRequest() {
	defer gock.Off()
	emailData := getEmailDetails()
	IncludeBaseTemplate := true

	expectedRequest := models.EmailRequest{
		From:    emailData.From,
		To:      emailData.To,
		Subject: emailData.Subject,
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  getEncodeMailContentWithBaseTemplate(),
		},
		IncludeBaseTemplate: IncludeBaseTemplate,
	}
	request, _ := json.Marshal(expectedRequest)

	gock.New(suite.gatewayURL).MatchType("json").JSON(request).Reply(200)

	err := suite.emailUtil.Send(emailData, IncludeBaseTemplate)

	assert.Nil(suite.T(), err)
}

func (suite *EmailUtilTest) TestShouldSendEmailWithAttachment_WhenValidRequest() {
	defer gock.Off()
	emailData := getEmailWithAttachmentDetails()
	IncludeBaseTemplate := true

	expectedRequest := models.EmailRequest{
		From:    emailData.From,
		To:      emailData.To,
		Subject: emailData.Subject,
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  getEncodeMailContentWithBaseTemplate(),
		},
		Attachments: []models.Attachment{
			{
				FileName: "advice.pdf",
				Data:     "PHRlc3QgY29udGVudD4=",
			},
		},
		IncludeBaseTemplate: IncludeBaseTemplate,
	}
	request, _ := json.Marshal(expectedRequest)

	gock.New(suite.gatewayURL).MatchType("json").JSON(request).Reply(200)

	err := suite.emailUtil.Send(emailData, IncludeBaseTemplate)

	assert.Nil(suite.T(), err)
}

func (suite *EmailUtilTest) TestReturnError_OnAPIFailure() {
	defer gock.Off()
	emailData := getEmailDetails()
	gock.New(suite.gatewayURL).MatchType("json").Reply(500)

	err := suite.emailUtil.Send(emailData, false)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), &golaerror.Error{
		ErrorCode:    "ERR_INTERNAL_SERVER_ERROR",
		ErrorMessage: "Error while sending Email",
	}, err)
}

func (suite *EmailUtilTest) TestSendEmailWithContext_ShouldSendEmail_WhenValidRequest() {
	defer gock.Off()
	emailData := getEmailDetails()
	IncludeBaseTemplate := true

	expectedRequest := models.EmailRequest{
		From:    emailData.From,
		To:      emailData.To,
		Subject: emailData.Subject,
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  getEncodeMailContentWithBaseTemplate(),
		},
		IncludeBaseTemplate: IncludeBaseTemplate,
	}
	request, _ := json.Marshal(expectedRequest)

	gock.New(suite.gatewayURL).MatchType("json").JSON(request).Reply(200)

	err := suite.emailUtil.SendWithContext(suite.goContext, emailData, IncludeBaseTemplate)

	assert.Nil(suite.T(), err)
}

func (suite *EmailUtilTest) TestSendEmailWithContext_ShouldSendEmailWithAttachment_WhenValidRequest() {
	defer gock.Off()
	emailData := getEmailWithAttachmentDetails()
	IncludeBaseTemplate := true

	expectedRequest := models.EmailRequest{
		From:    emailData.From,
		To:      emailData.To,
		Subject: emailData.Subject,
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  getEncodeMailContentWithBaseTemplate(),
		},
		Attachments: []models.Attachment{
			{
				FileName: "advice.pdf",
				Data:     "PHRlc3QgY29udGVudD4=",
			},
		},
		IncludeBaseTemplate: IncludeBaseTemplate,
	}
	request, _ := json.Marshal(expectedRequest)

	gock.New(suite.gatewayURL).MatchType("json").JSON(request).Reply(200)

	err := suite.emailUtil.SendWithContext(suite.goContext, emailData, IncludeBaseTemplate)

	assert.Nil(suite.T(), err)
}

func (suite *EmailUtilTest) TestSendEmailWithContext_ReturnError_OnAPIFailure() {
	defer gock.Off()
	emailData := getEmailDetails()
	gock.New(suite.gatewayURL).MatchType("json").Reply(500)

	err := suite.emailUtil.SendWithContext(suite.goContext, emailData, false)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), &golaerror.Error{
		ErrorCode:    "ERR_INTERNAL_SERVER_ERROR",
		ErrorMessage: "Error while sending Email",
	}, err)
}

func getEncodeMailContentWithBaseTemplate() string {
	return "PHA+VG9kYXkgaXMgYSBzdW5ueSBkYXk8L3A+"
}

func getEmailDetails() models.EmailDetails {
	sampleHTML := `<p>Today is a sunny day</p>`
	fromID := "no-reply@gola.xyz"
	toIDs := []string{"john@gmail.com"}
	subject := "Test Subject"
	emailData := models.EmailDetails{
		From:    fromID,
		To:      toIDs,
		Subject: subject,
		Content: sampleHTML,
	}
	return emailData
}

func getEmailWithAttachmentDetails() models.EmailDetails {
	emailData := getEmailDetails()
	emailData.Attachments = []models.EmailAttachment{
		{
			FileName: "advice.pdf",
			Content:  []byte("<test content>"),
		},
	}
	return emailData
}
