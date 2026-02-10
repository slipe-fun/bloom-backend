package google

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuthService struct {
	Config *oauth2.Config
}

func NewGoogleAuthService(clientID, clientSecret, redirectURL string) *GoogleAuthService {
	return &GoogleAuthService{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *GoogleAuthService) GetAuthURL(state string) string {
	return s.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *GoogleAuthService) ExchangeCode(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	token, err := s.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("cannot exchange code: %w", err)
	}
	return token, nil
}

func (s *GoogleAuthService) GetUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := s.Config.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed decode user info: %w", err)
	}

	return info, nil
}
