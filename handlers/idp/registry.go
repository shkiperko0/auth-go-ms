package idp

import (
	"net/http"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	Id          string
	Username    string
	DisplayName string
	Email       string
	AvatarUrl   string
}

type IHandlerProvider interface {
	SetHttpClient(client *http.Client)
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
	LoadConfig(clientId string, clientSecret string, redirectUrl string)
}

type IHandlerProviderRegistry interface {
	Add(method string, handler IHandlerProvider)
	Get(method string) IHandlerProvider
}

type HandlerProviderRegistry struct {
	items map[string]IHandlerProvider
}

func (h HandlerProviderRegistry) Add(method string, handler IHandlerProvider) {
	h.items[method] = handler
}

func (h HandlerProviderRegistry) Get(method string) IHandlerProvider {
	return h.items[method]
}

func NewHandlerProviderRegistry() IHandlerProviderRegistry {
	return &HandlerProviderRegistry{
		items: make(map[string]IHandlerProvider),
	}
}
