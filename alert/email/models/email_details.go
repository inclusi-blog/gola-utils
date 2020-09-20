package models

type EmailDetails struct {
	From        string
	To          []string
	Subject     string
	Content     string
	Attachments []EmailAttachment
}

type EmailAttachment struct {
	FileName string
	Content  []byte
}
