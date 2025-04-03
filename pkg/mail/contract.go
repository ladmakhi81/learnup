package mail

type Mail interface {
	SendPlain(dto SendMailReq) error
	SendTemplate(dto SendTemplateMailReq) error
}
