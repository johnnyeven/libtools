package g7

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/johnnyeven/libtools/conf"
	"github.com/johnnyeven/libtools/conf/presets"
	"github.com/johnnyeven/libtools/courier/client"
)

type OAuth struct {
	ClientID     string           `conf:"env"`
	ClientSecret presets.Password `conf:"env"`
	RedirectURL  string           `conf:"env"`
	Scopes       []string
	client.Client
}

func (o OAuth) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Host": "oauth.chinawayltd.com",
	}
}

func (o OAuth) MarshalDefaults(v interface{}) {
	if oauth, ok := v.(*OAuth); ok {
		if oauth.Host == "" {
			oauth.Host = "oauth.chinawayltd.com"
		}
	}
}

func (o *OAuth) authURL(baseURL, state string) string {
	var buf bytes.Buffer
	buf.WriteString(baseURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {o.ClientID},
		"redirect_uri":  {o.RedirectURL},
		"scope":         {strings.Join(o.Scopes, " ")},
		"state":         {state},
	}
	if strings.Contains(baseURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

func (o *OAuth) AuthCodeURL(state string) string {
	return o.authURL(o.Client.GetBaseURL("http")+"/oauth/authorize", state)
}

func (o *OAuth) LogoutURL(state string) string {
	return o.authURL(o.Client.GetBaseURL("http")+"/oauth/v0/logout", state)
}

type TokenRequest struct {
	// `Basic ${base64(client_id+":"+"client_secret")}`
	Authorization string `name:"Authorization" in:"header"`
	// required when grant_type=authorization_code
	Code string `name:"code" in:"formData"`
	// required when grant_type=refresh_token
	RefreshToken string `name:"refresh_token" in:"formData"`
	// authorization_code | password | client_credentials | refresh_token
	GrantType string `name:"grant_type" in:"formData"`
	// required when grant_type=password
	Username string `name:"username" in:"formData"`
	// required when grant_type=password
	Password string `name:"password" in:"formData"`
}

func (o *OAuth) Exchange(code string) (token *Token, err error) {
	token = &Token{}
	req := TokenRequest{
		GrantType:     "authorization_code",
		Code:          code,
		Authorization: "Basic " + basicAuth(o.ClientID, o.ClientSecret.String()),
	}
	err = o.Request("oauth.Token", "POST", "/oauth/token", req).
		Do().
		Into(token)
	return
}

type CheckTokenRequest struct {
	// `Basic ${base64(client_id+":"+"client_secret")}`
	Authorization string `name:"Authorization" in:"header"`
	// AccessToken
	AccessToken string `name:"access_token" in:"query"`
}

func (o OAuth) Validate(accessToken string) (token *Token, err error) {
	token = &Token{}
	req := CheckTokenRequest{
		AccessToken:   accessToken,
		Authorization: "Basic " + basicAuth(o.ClientID, o.ClientSecret.String()),
	}
	err = o.Request("oauth.CheckToken", "GET", "/oauth/token", req).
		Do().
		Into(token)
	return
}

type Token struct {
	TokenType   string `json:"token_type" default:"bearer"`
	AccessToken string `json:"access_token"`
	// 单位 s
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User
}

type User struct {
	ID uint64 `json:"id" redis:"id"`
	// ldap user
	Username       string `json:"username,omitempty"`
	Name           string `json:"name,omitempty"`
	OrgName        string `json:"org_name,omitempty"`
	GithubUsername string `json:"github_username,omitempty"`
	DingtalkID     string `json:"dingtalk_id,omitempty"`
	DingtalkNick   string `json:"dingtalk_nick,omitempty"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
