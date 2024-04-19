package mailQueue

import (
	"bytes"
	"crypto/tls"
	"errors"
	"gopkg.in/gomail.v2"
	"html/template"
	"neiro-api/config"
	"neiro-api/internal/services/mail"
)

func mailQueueHandler(data mail.Data) error {
	cfg := config.GetConfig().MailConfig
	msg := gomail.NewMessage()
	from := cfg.FromMailer
	if data.From != "" {
		from = data.From
	}
	msg.SetHeader("From", from)
	msg.SetHeader("To", data.To)
	msg.SetHeader("Subject", data.Subject)

	t, err := template.ParseFiles("assets/email/template.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer

	type TemplateData struct {
		PreHeader string
		Messages  []string
	}
	templateData := TemplateData{
		PreHeader: data.Header,
		Messages:  data.Lines,
	}
	if err := t.Execute(&tpl, templateData); err != nil {
		return err
	}
	msg.SetBody("text/html", tpl.String())

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Login, cfg.Password)
	if cfg.TLS {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if err = d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func TaskHandlerWrapper(data interface{}) error {
	mailData, ok := data.(mail.Data)
	if !ok {
		return errors.New("не верный тип данных, ожидался mail.Data")
	}
	return mailQueueHandler(mailData)
}
