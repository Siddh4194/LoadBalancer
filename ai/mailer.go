package ai

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	gmailSMTPServer = "smtp.gmail.com"
	gmailSMTPPort   = 587
)

func buildSMTPMessage(fromEmail, toEmail, subject, htmlBody string) string {
	headers := map[string]string{
		"From":         fromEmail,
		"To":           toEmail,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	var msg strings.Builder
	for key, value := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	return msg.String()
}

// SendHTMLReportSMTP sends an HTML report through Gmail SMTP using the sender email and Gmail app password.
//
// The Gmail account must have 2-Step Verification enabled and an app password created.
func SendHTMLReportSMTP(subject string, htmlBody string) error {
	fromEmail := "siddh4194@gmail.com"
	appPassword := "zanhzpgygbgwkkqd" // Replace with your Gmail app password
	toEmail := "siddhantkadam.dev@gmail.com"
	auth := smtp.PlainAuth("", fromEmail, appPassword, gmailSMTPServer)

	message := buildSMTPMessage(fromEmail, toEmail, "The report of the actions taken by the load balancer's ai", htmlBody)
	addr := fmt.Sprintf("%s:%d", gmailSMTPServer, gmailSMTPPort)
	return smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, []byte(message))
}
