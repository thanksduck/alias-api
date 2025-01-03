package utils

import (
	"net/smtp"
	"os"
)

func SendEmail(to string, subject string, htmlBody string, textBody string) error {
	from := os.Getenv("EMAIL_FROM")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	boundary := "boundary123"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: multipart/alternative; boundary=" + boundary + "\n\n" +
		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=utf-8\n\n" +
		textBody + "\n\n" +
		"--" + boundary + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		htmlBody + "\n\n" +
		"--" + boundary + "--"

	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

/*

func SendEmailAsync(to, subject, body string) chan error {
	resultChan := make(chan error, 1)
	go func() {
		from := os.Getenv("EMAIL_FROM")
		username := os.Getenv("EMAIL_USERNAME")
		password := os.Getenv("EMAIL_PASSWORD")
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")
		msg := "From: " + from + "\n" +
			"To: " + to + "\n" +
			"Subject: " + subject + "\n\n" +
			body
		auth := smtp.PlainAuth("", username, password, smtpHost)
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
		resultChan <- err
	}()
	return resultChan
}
*/
