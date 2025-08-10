package utils

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauth struct {
	Config *oauth2.Config
	UrlAPI string
}

func CreateGoogleUtil() *GoogleOauth {
	redirectUrl := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")

	return &GoogleOauth{
		Config: &oauth2.Config{
			RedirectURL:  redirectUrl,
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		UrlAPI: "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
	}
}
