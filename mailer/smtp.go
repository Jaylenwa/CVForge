package mailer

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"openresume/config"
)

func SendVerificationCode(cfg config.Config, toEmail, code string) error {
	subject := "Your OpenResume verification code"
	body := fmt.Sprintf("<div style=\"font-family:Arial,sans-serif;font-size:14px;color:#333\">"+
		"<p>Your verification code is:</p>"+
		"<p style=\"font-size:24px;font-weight:bold;letter-spacing:4px\">%s</p>"+
		"<p>The code expires in 10 minutes.</p>"+
		"<p>If you did not request this, please ignore.</p>"+
		"</div>", code)
	return sendSMTP(cfg, toEmail, subject, body)
}

func sendSMTP(cfg config.Config, toEmail, subject, htmlBody string) error {
	host := cfg.SMTPHost
	portStr := cfg.SMTPPort
	user := cfg.SMTPUsername
	pass := cfg.SMTPPassword
	fromName := cfg.SMTPFromName
	if host == "" || portStr == "" || user == "" || pass == "" {
		return fmt.Errorf("smtp config missing")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}
	from := mail.Address{Name: fromName, Address: user}
	to := mail.Address{Address: toEmail}
	headers := map[string]string{
		"From":         from.String(),
		"To":           to.String(),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}
	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(k)
		msgBuilder.WriteString(": ")
		msgBuilder.WriteString(v)
		msgBuilder.WriteString("\r\n")
	}
	msgBuilder.WriteString("\r\n")
	msgBuilder.WriteString(htmlBody)
	addr := fmt.Sprintf("%s:%d", host, port)
	tlsCfg := &tls.Config{ServerName: host}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()
	if ok, _ := client.Extension("AUTH"); ok {
		auth := smtp.PlainAuth("", user, pass, host)
		if err = client.Auth(auth); err != nil {
			return err
		}
	}
	if err = client.Mail(from.Address); err != nil {
		return err
	}
	if err = client.Rcpt(to.Address); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write([]byte(msgBuilder.String())); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return nil
}
