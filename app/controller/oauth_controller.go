package controller

import (
	"net/http"
)

type OAuthController struct{}

func NewOAuthController() OAuthController {
	return OAuthController{}
}

func (h OAuthController) StartOAuth(c Context) error {
	return c.String(http.StatusOK, "start")
}