package email

// mockgen -source=alert/email/email_util.go -destination=mocks/mock_email_util.go -package=mocks
import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/gola-glitch/gola-utils/alert/email/models"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/http/request"
	"github.com/gola-glitch/gola-utils/logging"
	"go.opencensus.io/plugin/ochttp"
)

type Util interface {
	Send(emailDetails models.EmailDetails, includeBaseTemplate bool) *golaerror.Error
	SendWithContext(context context.Context, emailDetails models.EmailDetails, includeBaseTemplate bool) *golaerror.Error
}

type emailUtil struct {
	gatewayURL string
}

func NewEmailUtil(gatewayURL string) Util {
	return emailUtil{gatewayURL}
}

func (util emailUtil) Send(emailDetails models.EmailDetails, includeBaseTemplate bool) *golaerror.Error {
	contentByteArray := []byte(emailDetails.Content)
	encodedMailContent := base64.StdEncoding.EncodeToString(contentByteArray)

	emailRequest := mapToEmailRequest(emailDetails, encodedMailContent, includeBaseTemplate)
	return sendEmail(context.TODO(), emailRequest, util.gatewayURL)
}

func (util emailUtil) SendWithContext(context context.Context, emailDetails models.EmailDetails, includeBaseTemplate bool) *golaerror.Error {
	contentByteArray := []byte(emailDetails.Content)
	encodedMailContent := base64.StdEncoding.EncodeToString(contentByteArray)

	emailRequest := mapToEmailRequest(emailDetails, encodedMailContent, includeBaseTemplate)
	return sendEmailWithLoggingContext(context, emailRequest, util.gatewayURL)
}

func mapToEmailRequest(emailDetails models.EmailDetails, encodedContent string, includeBaseTemplate bool) models.EmailRequest {
	emailRequest := models.EmailRequest{
		From:    emailDetails.From,
		To:      emailDetails.To,
		Subject: emailDetails.Subject,
		Body: models.MessageBody{
			MimeType: "text/html",
			Content:  encodedContent,
		},
		IncludeBaseTemplate: includeBaseTemplate,
	}
	if emailDetails.Attachments != nil {
		encodedAttachments := make([]models.Attachment, len(emailDetails.Attachments))
		for i, attachment := range emailDetails.Attachments {
			encodedAttachments[i] = models.Attachment{
				FileName: attachment.FileName,
				Data:     base64.StdEncoding.EncodeToString(attachment.Content),
			}
		}
		emailRequest.Attachments = encodedAttachments
	}
	return emailRequest
}

func sendEmail(context context.Context, emailRequest models.EmailRequest, gatewayURL string) *golaerror.Error {
	logger := logging.GetLogger(context)
	httpClient := &http.Client{Transport: &ochttp.Transport{}}
	responseError := request.NewHttpRequestBuilder(httpClient).
		NewRequest().
		WithContext(context).
		WithJSONBody(emailRequest).
		Post(gatewayURL)
	if responseError != nil {
		logger.WithField("error", responseError).Error("Error while sending Email")
		return &golaerror.Error{ErrorCode: "ERR_INTERNAL_SERVER_ERROR", ErrorMessage: "Error while sending Email"}
	}
	return nil
}

func sendEmailWithLoggingContext(context context.Context, emailRequest models.EmailRequest, gatewayURL string) *golaerror.Error {
	logger := logging.GetLogger(context)
	httpClient := &http.Client{Transport: &ochttp.Transport{}}
	responseError := request.NewHttpRequestBuilder(httpClient).
		NewRequest().
		WithContext(context).
		WithJSONBody(emailRequest).
		Post(gatewayURL)
	if responseError != nil {
		logger.Error("Error while sending Email :", responseError)
		return &golaerror.Error{ErrorCode: "ERR_INTERNAL_SERVER_ERROR", ErrorMessage: "Error while sending Email"}
	}
	return nil
}
