package models

type EmailRequest struct {
	From                string       `json:"from"`
	To                  []string     `json:"to"`
	Subject             string       `json:"subject"`
	Body                MessageBody  `json:"message_body"`
	Attachments         []Attachment `json:"attachments"`
	IncludeBaseTemplate bool         `json:"include_base_template"`
}
type MessageBody struct {
	MimeType string `json:"mime_type"`
	Content  string `json:"base64_encoded_content"`
}

type Attachment struct {
	FileName string `json:"file_name"`
	Data     string `json:"base64_encoded_data"`
}
