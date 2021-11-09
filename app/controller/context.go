package controller

type Context interface {
	String(code int, s string) error
	HTML(code int, html string) error
	Render(code int, name string, data interface{}) error
}