package mask_util

type MaskFieldKeyType string

const (
	maskChar     = "*"
	fallbackText = "*NA*"
)

const (
	Email         MaskFieldKeyType = "email"
)
