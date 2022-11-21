package idp

import (
	"net/http"

	"golang.org/x/oauth2"
)

type IdProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

type IdpProvider interface {
	Init() *IdProvider
}

func (h *IdProvider) Init() *IdProvider {
	idp := &IdProvider{}
	return idp
}

func (h *IdProvider) SetConfig(config *oauth2.Config) {
	h.Config = config
}

func (h *IdProvider) SetHttpClient(client *http.Client) {
	h.Client = client
}

func NewIdpProvider() IdpProvider {
	return &IdProvider{}
}
