package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

type (
	gmail struct {
		cfg     *Parameters
		headers textproto.MIMEHeader
	}
)

func New(cfg *Parameters) IEmail {
	h := map[string]string{
		"Date":                        time.Now().Format(layoutDateEmailSend),
		"Subject":                     cfg.Subject,
		"Sender":                      cfg.From,
		"Return-Path":                 cfg.From,
		"Reply-To":                    cfg.From,
		"From":                        cfg.From,
		"To":                          cfg.To,
		"Disposition-Notification-To": cfg.From,
		"Return-Receipt-To":           cfg.From,
		"MIME-Version":                fmt.Sprintf("1.0;\nContent-Type: %s; charset: %s;", ContentType, Charset),
	}

	headers := make(textproto.MIMEHeader, 0)

	for k, v := range h {
		headers[k] = []string{v}
	}

	return &gmail{
		cfg:     cfg,
		headers: headers,
	}
}

func (g *gmail) Send() error {
	var message, mime string

	for _, v := range PriorityHeaders {
		mime += fmt.Sprintf(MessageFormat, v.Header, v.Value)
	}

	for k, v := range g.headers {
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

func (g *gmail) Properties() *Properties {
	emailHeaders := make(map[string]string, 0)

	for k, v := range g.headers {
		emailHeaders[k] = strings.Join(v, ",")
	}

	return &Properties{
		From:    g.cfg.From,
		To:      g.cfg.To,
		Headers: emailHeaders,
	}
}
