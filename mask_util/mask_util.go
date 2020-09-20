package mask_util

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"strings"
)

func MaskEmail(ctx context.Context, value string) string {
	return validateAndMaskField(ctx, Email, value)
}

func validateAndMaskField(ctx context.Context, key MaskFieldKeyType, value string) string {
	log := logging.GetLogger(ctx).WithField("class", "mask_util").WithField("method", "validateAndMaskField")
	if strings.TrimSpace(value) == "" {
		log.Error("cannot mask empty string, so return fallback text")
		return fallbackText
	}

	switch key {
	case Email:
		return maskEmail(ctx, value)
	default:
		log.Warnf("There is no key found with the name of %v, so return fallback text", key)
		return fallbackText
	}
}

func maskString(ctx context.Context, targetString string, skipFromFirst, skipFromLast int) string {
	log := logging.GetLogger(ctx).WithField("class", "mask_util").WithField("method", "maskString")
	sumOfMaskCharLen := skipFromFirst + skipFromLast
	if (sumOfMaskCharLen) > len(targetString) {
		log.Error("error while masking string")
		return fallbackText
	}
	masked := targetString[:skipFromFirst] + strings.Repeat(maskChar, len(targetString)-(sumOfMaskCharLen)) + targetString[len(targetString)-skipFromLast:]
	return masked
}

func maskEmail(ctx context.Context, email string) string {
	array := strings.Split(email, "@")
	domainArray := strings.Split(array[1], ".")
	username := maskString(ctx, array[0], 1, 1)
	mailSuffix := maskString(ctx, domainArray[0], 1, 1)
	maskedEmail := username + "@" + mailSuffix + "." + strings.Join(domainArray[1:], ".")
	return maskedEmail
}
