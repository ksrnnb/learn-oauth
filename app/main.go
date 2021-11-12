package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/ksrnnb/learn-oauth/app/middleware"
	"github.com/ksrnnb/learn-oauth/app/route"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	middleware.SetMiddleware(e)
	route.SetRoute(e)

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err)
		return
	}

	e.Logger.Fatal(e.Start(":3000"))
}
