package notify

import (
	"log"
	"net/smtp"
)

type EmailClient struct {
	EmailAccount	string
	AppPassword		string
	SmtpServer		string
	SmtpPort		string
}

func (email *EmailClient) SendEmail(destinationAddr string, subjectLine string, emailBody string) error {

	emailSrc := "From: " + email.EmailAccount + "\n"
	emailDest := "To: " + destinationAddr + "\r\n"
	subject := "Subject: " + subjectLine + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	rawMessage := []byte(emailSrc + emailDest + subject + mime + "\n" + emailBody)
	// encodedMessage := base64.URLEncoding.EncodeToString(rawMessage)

	err := smtp.SendMail(email.SmtpServer + ":" + email.SmtpPort, 
						smtp.PlainAuth("", email.EmailAccount, email.AppPassword, email.SmtpServer), 
						email.EmailAccount, []string{destinationAddr}, 
						rawMessage)

	if err != nil {
		log.Printf("smtp error: %s", err)
	}

	return err
}
