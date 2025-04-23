package smtp

import (
	"errors"
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/dtos"
	"github.com/ladmakhi81/learnup/shared/utils"
	"net/smtp"
)

type SmtpMailSvc struct {
	config *dtos.EnvConfig
}

func NewSmtpMailSvc(config *dtos.EnvConfig) *SmtpMailSvc {
	return &SmtpMailSvc{
		config: config,
	}
}

func (svc SmtpMailSvc) SendPlain(dto dtos.SendMailReq) error {
	receivers := []string{dto.Recipient}
	body := fmt.Sprintf(`To: %s
Subject: %s

%s
`, dto.Recipient, dto.Subject, dto.Body)
	addr := svc.getAddr()
	auth := svc.getSmtpAuth()
	mailUsername := svc.config.Smtp.Username
	err := smtp.SendMail(addr, auth, mailUsername, receivers, []byte(body))
	if err != nil {
		return errors.New("Error: happen in sending plain email")
	}
	return nil
}

func (svc SmtpMailSvc) SendTemplate(dto dtos.SendTemplateMailReq) error {
	receivers := []string{dto.Recipient}
	addr := svc.getAddr()
	auth := svc.getSmtpAuth()
	parsedTemplate, parsedErr := utils.ParseTemplate(fmt.Sprintf("mail/%s", dto.Template), dto.TemplateData)
	if parsedErr != nil {
		return errors.New("Error: happen in finding template")
	}
	mailUsername := svc.config.Smtp.Username
	body := fmt.Sprintf(`MIME-Version: 1.0
Content-Type: text/html; charset: utf-8;
From: %s
To: %s
Subject: %s

%s`, mailUsername, dto.Recipient, dto.Subject, parsedTemplate)
	err := smtp.SendMail(addr, auth, mailUsername, receivers, []byte(body))
	if err != nil {
		return errors.New("Error: happen in sending template email")
	}
	return nil
}

func (svc SmtpMailSvc) getSmtpAuth() smtp.Auth {
	smtpConfig := svc.config.Smtp

	return smtp.PlainAuth(
		"",
		smtpConfig.Username,
		smtpConfig.Password,
		smtpConfig.Host,
	)
}

func (svc SmtpMailSvc) getAddr() string {
	return fmt.Sprintf("%s:%s", svc.config.Smtp.Host, svc.config.Smtp.Port)
}
