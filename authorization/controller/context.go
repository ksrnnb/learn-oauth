package controller

type Context interface {
	String(code int, s string) error
	HTML(code int, html string) error
	Render(code int, name string, data interface{}) error
	Get(key string) interface{}
	FormValue(name string) string
	Redirect(code int, url string) error
}