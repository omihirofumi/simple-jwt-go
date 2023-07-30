package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	YOUR_ID       = "hirofumi"
	YOUR_PASSWORD = "mypass"
)

type userInfo struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)

	r := e.Group("/restricted")
	{
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(jwtClaims)
			},
			SigningKey: []byte(KEY_NAME)}
		r.Use(echojwt.WithConfig(config))
		r.GET("/hello", helloFunc)
	}
	e.Logger.Fatal(e.Start(":50051"))
}
