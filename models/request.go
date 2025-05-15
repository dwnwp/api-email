package models

type MailerRequest struct {
	From    string `json:"From"`
	To      string `json:"To"`
	Subject string `json:"Subject"`
	BodySubject    string `json:"BodySubject"`
	BodyContent    string `json:"BodyContent"`
}
