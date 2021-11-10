package session

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const sessionName = "session"

// セッションの保存
func Save(key interface{}, value interface{}, c echo.Context) error {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return err
	}

	sess.Values[key] = value
	return sess.Save(c.Request(), c.Response())
}

func Get(key interface{}, c echo.Context) (interface{}, error) {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return nil, err
	}

	value, ok := sess.Values[key]

	if !ok {
		return nil, errors.New("session " + key.(string) + " doesn't exist")
	}

	return value, nil

}

// セッション破棄
func Destroy(c echo.Context) error {
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{MaxAge: -1, Path: "/"}
    return sess.Save(c.Request(), c.Response())
}