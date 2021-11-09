package controller

import (
	"net/http"
)

type HomeController struct{}

func NewHomeController() HomeController {
	return HomeController{}
}

func (h HomeController) Home(c Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"csrf": c.Get("csrf"),
	})
}