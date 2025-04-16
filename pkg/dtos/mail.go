package dtos

type SendMailReq struct {
	Recipient string
	Subject   string
	Body      string
}

func NewSendMailReq(recipient, subject, body string) SendMailReq {
	return SendMailReq{
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
	}
}

type SendTemplateMailReq struct {
	Recipient    string
	Subject      string
	Template     string
	TemplateData any
}

func NewSendTemplateMailReq(recipient, subject, template string, templateData any) SendTemplateMailReq {
	return SendTemplateMailReq{
		Recipient:    recipient,
		Subject:      subject,
		Template:     template,
		TemplateData: templateData,
	}
}
