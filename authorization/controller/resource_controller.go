package controller

import (
	"net/http"
	"strings"

	"github.com/ksrnnb/learn-oauth/authorization/resource"

	"github.com/labstack/echo/v4"
)

type ResourceController struct{}

func NewResourceController() ResourceController {
	return ResourceController{}
}

type ResourceResponse struct {
	UserId     int    `json:"userId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	PictureUrl string `json:"pictureUrl"`
}

// リソース情報を取得する
func (controller ResourceController) Resource(c echo.Context) error {
	authorizationHeader := c.Request().Header.Get("Authorization")
	splitted := strings.Split(authorizationHeader, "Bearer ")

	if len(splitted) != 2 {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "authorization header is invalid")
	}

	tokenString := splitted[1]
	token, err := resource.FindAccessTokenFromToken(tokenString)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "bearer token is invalid")
	}

	if token.Expired() {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "bearer token has expired")
	}

	user, err := token.FindUser()

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "user is not found")
	}

	res := &ResourceResponse{
		UserId:     user.Id,
		Name:       user.Name,
		Email:      user.Email,
		PictureUrl: user.PictureUrl,
	}

	return c.JSON(http.StatusOK, res)
}
