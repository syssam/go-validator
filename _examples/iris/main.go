package main

import (
	"net/http"

	"github.com/kataras/iris"
	validator "github.com/syssam/go-validator"
)

type User struct {
	Name  string `json:"name" valid:"required"`
	Email string `json:"email" valid:"required,email"`
}

var validate *validator.Validator

func main() {
	app := iris.New()
	validate = validator.New()

	app.Post("/user", func(c iris.Context) {
		var user User
		if err := c.ReadJSON(&user); err != nil {
			// Handle error.
		}

		if err := validator.ValidateStruct(user); err != nil {
			c.StatusCode(http.StatusBadRequest)
			c.JSON(err)
			return
		}

		c.StatusCode(http.StatusOK)
		c.JSON(user)
	})

	app.Run(iris.Addr(":8080"))
}
