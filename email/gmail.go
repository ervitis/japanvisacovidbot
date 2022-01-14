package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"time"
)

type (
	gmail struct {
		cfg      *Parameters
		headers  textproto.MIMEHeader
		encoding string
	}
)

func New(cfg *Parameters) IEmail {
	headers["Sender"] = cfg.From
	headers["Return-Path"] = cfg.From
	headers["Reply-To"] = cfg.From
	headers["From"] = cfg.From
	headers["To"] = cfg.To
	headers["Disposition-Notification-To"] = cfg.From

	return &gmail{
		cfg: cfg,
	}
}

func (g *gmail) Send() error {
	var message, mime string

	for _, v := range PriorityHeaders {
		mime += fmt.Sprintf(MessageFormat, v.Header, v.Value)
	}

	for k, v := range headers {
		mime += fmt.Sprintf(MessageFormat, k, v)
	}

	loc, err := time.LoadLocation(locationEmail)
	if err != nil {
		return err
	}

	now := time.Now().In(loc)

	s := "Buenos dias"
	if now.Hour() <= 12 && now.Hour() <= 19 {
		s = "Buenas tardes"
	} else if now.Hour() > 19 {
		s = "Buenas noches"
	}
	body := fmt.Sprintf(MessageBody, s)

	message += fmt.Sprintf("%s: %s\n%s\n%s", "Subject", g.cfg.Subject, mime, body)

	host, _, _ := net.SplitHostPort(fmt.Sprintf(`%s:%s`, g.cfg.SmtpServer, g.cfg.SmtpPort))
	auth := smtp.PlainAuth("", g.cfg.Username, g.cfg.Password, host)
	tlsConfig := &tls.Config{InsecureSkipVerify: false, ServerName: host}

	svc, err := smtp.Dial(fmt.Sprintf(`%s:%s`, g.cfg.SmtpServer, g.cfg.SmtpPort))
	if err != nil {
		return err
	}

	if err := svc.StartTLS(tlsConfig); err != nil {
		return err
	}

	if err := svc.Auth(auth); err != nil {
		return err
	}

	if err := svc.Mail(g.cfg.From); err != nil {
		return err
	}

	if err := svc.Rcpt(g.cfg.To); err != nil {
		return err
	}

	w, err := svc.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return svc.Quit()
}
