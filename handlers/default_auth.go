package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/shkiperko0/auth-go-ms/handlers/idp"

	"golang.org/x/oauth2"
)

type IDefaultAuthHanlder interface {
}

type DefaultAuthHanlder struct {
	Provider *idp.IdProvider
}

func NewDefaultAuthHandler(Provider *idp.IdProvider) idp.IHandlerProvider {
	return &DefaultAuthHanlder{
		Provider: Provider,
	}
}

func (i *DefaultAuthHanlder) SetHttpClient(client *http.Client) {
	i.Provider.SetHttpClient(client)
}

func (i *DefaultAuthHanlder) LoadConfig(clientId string, clientSecret string, redirectUrl string) {
	var endpoint = oauth2.Endpoint{}

	var config = &oauth2.Config{
		Scopes:       []string{"profile", "email"},
		Endpoint:     endpoint,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}

	i.Provider.SetConfig(config)
}

func (i *DefaultAuthHanlder) GetToken(code string) (*oauth2.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, i.Provider.Client)
	return i.Provider.Config.Exchange(ctx, code)
}

type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (i *DefaultAuthHanlder) GetUserInfo(token *oauth2.Token) (*idp.UserInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=%s", token.AccessToken)
	resp, err := i.Provider.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUserInfo GoogleUserInfo
	err = json.Unmarshal(body, &googleUserInfo)
	if err != nil {
		return nil, err
	}

	if googleUserInfo.Email == "" {
		return nil, errors.New("google email is empty")
	}

	userInfo := idp.UserInfo{
		Id:          googleUserInfo.Id,
		Username:    googleUserInfo.Email,
		DisplayName: googleUserInfo.Name,
		Email:       googleUserInfo.Email,
		AvatarUrl:   googleUserInfo.Picture,
	}

	return &userInfo, nil
}
