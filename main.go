package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
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

	e.POST("/login", login)

	e.Logger.Fatal(e.Start(":50051"))
}

func login(c echo.Context) error {
	var loginInfo userInfo

	if err := c.Bind(&loginInfo); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(loginInfo); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !isValidUser(loginInfo.Username, loginInfo.Password) {
		return c.String(http.StatusUnauthorized, "")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("SECRET_KEY"))
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})

}

func isValidUser(username, password string) bool {
	return username == YOUR_ID && password == YOUR_PASSWORD
}
