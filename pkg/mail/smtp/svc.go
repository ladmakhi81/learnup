package smtp

import (
	"errors"
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/mail"
	"github.com/ladmakhi81/learnup/utils"
	"net/smtp"
)

// TODO: remove this struct when config file wrote (this is temporary and this struct must removed)
type SmtpConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type SmtpMailSvc struct {
}

func NewSmtpMailSvc() *SmtpMailSvc {
	return &SmtpMailSvc{}
}

func (svc SmtpMailSvc) SendPlain(dto mail.SendMailReq) error {
	receivers := []string{dto.Recipient}
	body := fmt.Sprintf(`To: %s
Subject: %s

%s
`, dto.Recipient, dto.Subject, dto.Body)
	addr := svc.getAddr()
	auth := svc.getSmtpAuth()
	//TODO: update this code when config file write
	smtpConfig := svc.getSmtpConfig()
	mailUsername := smtpConfig.Username
	err := smtp.SendMail(addr, auth, mailUsername, receivers, []byte(body))
	if err != nil {
		return errors.New("Error: happen in sending plain email")
	}
	return nil
}

func (svc SmtpMailSvc) SendTemplate(dto mail.SendTemplateMailReq) error {
	receivers := []string{dto.Recipient}
	addr := svc.getAddr()
	auth := svc.getSmtpAuth()
	parsedTemplate, parsedErr := utils.ParseTemplate(fmt.Sprintf("mail/%s", dto.Template), dto.TemplateData)
	if parsedErr != nil {
		return errors.New("Error: happen in finding template")
	}
	//TODO: update this code when config file write
	smtpConfig := svc.getSmtpConfig()
	mailUsername := smtpConfig.Username
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
	//TODO: update this code when config file write
	smtpConfig := svc.getSmtpConfig()

	return smtp.PlainAuth(
		"",
		smtpConfig.Username,
		smtpConfig.Password,
		smtpConfig.Host,
	)
}

func (svc SmtpMailSvc) getAddr() string {
	//TODO: update this code when config file write
	smtpConfig := svc.getSmtpConfig()
	return fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)
}

func (svc SmtpMailSvc) getSmtpConfig() SmtpConfig {
	return SmtpConfig{
		Host:     "host",
		Port:     "port",
		Username: "username",
		Password: "password",
	}
}
