package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

const (
	KEY_NAME = "SECRET_KEY"
)

type jwtClaims struct {
	ID string `json:":id"`
	jwt.RegisteredClaims
}

func login(c echo.Context) error {
	var loginInfo userInfo

	if err := c.Bind(&loginInfo); err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}

	validate := validator.New()
	if err := validate.Struct(loginInfo); err != nil {
		log.Println(err)
		return echo.ErrBadRequest
	}

	if !isValidUser(loginInfo.Username, loginInfo.Password) {
		return echo.ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		loginInfo.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})

	tokenString, err := token.SignedString([]byte(KEY_NAME))
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func helloFunc(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtClaims)
	return c.String(http.StatusOK, claims.ID)
}
