package oauth2

import (
	"golang.org/x/oauth2"

	"github.com/johnnyeven/libtools/conf/presets"
)

type OAuthConfig struct {
	ClientID     string           `conf:"env"`
	ClientSecret presets.Password `conf:"env"`
	RedirectURL  string           `conf:"env"`
	Scopes       []string
	oauth2.Endpoint
}

func (o OAuthConfig) Init() {
	oauth2.RegisterBrokenAuthHeaderProvider(o.Endpoint.TokenURL)
}

func (o OAuthConfig) Get() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: string(o.ClientSecret),
		Endpoint:     o.Endpoint,
		RedirectURL:  o.RedirectURL,
		Scopes:       o.Scopes,
	}
}
