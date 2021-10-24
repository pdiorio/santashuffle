package notify

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailClient struct {
	ClientID          string
	ClientSecret      string
	AccessToken       string
	RefreshToken      string
	AuthorizationCode string
	Connection        *gmail.Service
}

func (g *GmailClient) Connect() error {
	if (g.ClientID == "") || (g.ClientSecret == "") {
		return errors.New("gmail connect: must set a clientID and clientSecret")
	}

	oauth2Config := oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://mail.google.com/"},
		RedirectURL:  "http://localhost",
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
	//defer cancel()
	ctx := context.Background()

	var tokenPtr *oauth2.Token

	if (g.RefreshToken != "") && (g.AccessToken != "") {
		tokenPtr = &oauth2.Token{
			AccessToken:  g.AccessToken,
			RefreshToken: g.RefreshToken,
			TokenType:    "Bearer",
			Expiry:       time.Now(),
		}
	} else if g.AuthorizationCode != "" {
		var err error
		tokenPtr, err = oauth2Config.Exchange(ctx, g.AuthorizationCode)
		if err != nil {
			return err
		}
	} else {
		return errors.New("gmail connect: must pre-set either (accessToken & refreshToken) or authorizationCode")
	}

	tokenSource := oauth2Config.TokenSource(ctx, tokenPtr)
	service, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return errors.New("gmail connect: error creating gmail service from settings")
	}

	g.Connection = service

	return nil
}

func (g *GmailClient) SendEmail(destinationAddr string, subjectLine string, emailBody string) error {
	if g.Connection == nil {
		err := g.Connect()
		if err != nil {
			return err
		}
	}

	emailDest := "To: " + destinationAddr + "\r\n"
	subject := "Subject: " + subjectLine + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	rawMessage := []byte(emailDest + subject + mime + "\n" + emailBody)
	encodedMessage := base64.URLEncoding.EncodeToString(rawMessage)

	message := gmail.Message{Raw: encodedMessage}

	_, err := g.Connection.Users.Messages.Send("me", &message).Do()
	return err
}
