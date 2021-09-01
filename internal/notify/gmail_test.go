package notify

import (
	"testing"
)

func TestGmailServiceConnectNoClientInfo(t *testing.T) {
	gSrv := GmailClient{}
	err := gSrv.Connect()

	if err == nil {
		t.Errorf("GmailService.Connect() should return error for empty clientID and clientSecret")
	}
}

func TestGmailServiceConnectNoTokenInfo(t *testing.T) {
	gSrv := GmailClient{ClientID: "foo", ClientSecret: "bar"}
	err := gSrv.Connect()

	if err == nil {
		t.Errorf("GmailService.Connect() should return error for empty accessToken+refreshToken or authorizationCode")
	}
}

func TestGmailServiceConnectNoAuthCode(t *testing.T) {
	gSrv := GmailClient{ClientID: "foo", ClientSecret: "bar", AuthorizationCode: "fake"}
	err := gSrv.Connect()

	if err == nil {
		t.Errorf("GmailService.Connect() should return error with failed attempt to use authorizationCode")
	}
}

func TestGmailServiceConnectCleanService(t *testing.T) {
	gSrv := GmailClient{ClientID: "foo", ClientSecret: "bar", AccessToken: "crypt", RefreshToken: "ography"}
	err := gSrv.Connect()

	if err != nil {
		t.Errorf("GmailService.Connect() error in creating gmail service with properly configured client")
	}
}
