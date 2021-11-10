package session

import (
	"errors"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// セッションの保存
func Save(key interface{}, value interface{}, c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Values[key] = value
	return sess.Save(c.Request(), c.Response())
}

func Get(key interface{}, c echo.Context) (interface{}, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}

	value, ok := sess.Values[key]

	if !ok {
		return nil, errors.New("session " + key.(string) + " doesn't exist")
	}

	return value, nil

}