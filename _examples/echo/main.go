package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Name  string `json:"name" valid:"required"`
	Email string `json:"email" valid:"required,email"`
}

func main() {
	e := echo.New()
	e.Validator = &DefaultValidator{}
	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return
		}
		if err = c.Validate(u); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
