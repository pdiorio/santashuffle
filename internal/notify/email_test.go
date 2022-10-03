package notify

import (
	"testing"
)

func TestEmailServiceNoInfo(t *testing.T) {
	eSrv := EmailClient{}
	err := eSrv.SendEmail("fake2@gmail.com", "Test Email", "This is a test.")

	if err == nil {
		t.Errorf("EmailService.SendEmail should return error when provided no connection information")
	}
}

func TestEmailServiceInvalidInfo(t *testing.T) {
	eSrv := EmailClient{EmailAccount: "fake@gmail.com", AppPassword: "abc123", SmtpServer: "smtp.gmail.com", SmtpPort: "587"}
	err := eSrv.SendEmail("fake2@gmail.com", "Test Email", "This is a test.")

	if err == nil {
		t.Errorf("EmailService.SendEmail should return error when provided invalid connection information")
	}
}
