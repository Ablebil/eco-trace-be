package oauth

import (
	"context"
	"strings"

	"github.com/Ablebil/eco-sample/config"
	"github.com/Ablebil/eco-sample/internal/domain/dto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type OAuthItf interface {
	GenerateLink(state string) (string, error)
	ExchangeToken(code string) (*oauth2.Token, error)
	GetProfile(token *oauth2.Token) (*dto.GoogleProfileResponse, error)
}

type OAuth struct {
	googleOAuthConfig *oauth2.Config
}

func NewOAuth(cfg *config.Config) OAuthItf {
	return &OAuth{
		googleOAuthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (o *OAuth) GenerateLink(state string) (string, error) {
	return o.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (o *OAuth) ExchangeToken(code string) (*oauth2.Token, error) {
	token, err := o.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return &oauth2.Token{}, err
	}

	return token, nil
}

func (o *OAuth) GetProfile(token *oauth2.Token) (*dto.GoogleProfileResponse, error) {
	client := o.googleOAuthConfig.Client(context.Background(), token)
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	data := &dto.GoogleProfileResponse{
		ID:       userInfo.Id,
		Email:    userInfo.Email,
		Username: strings.Split(userInfo.Email, "@")[0],
		Name:     userInfo.Name,
		Verified: *userInfo.VerifiedEmail,
	}

	return data, nil
}
