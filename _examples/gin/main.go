package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/syssam/go-validator"
	lang_en "github.com/syssam/go-validator/_examples/translations/lang/en"
	lang_zhCN "github.com/syssam/go-validator/_examples/translations/lang/zh_CN"
	validator_en "github.com/syssam/go-validator/lang/en"
	validator_zhCN "github.com/syssam/go-validator/lang/zh_CN"
)

// User contains user information
type User struct {
	FirstName string `valid:"required"`
	LastName  string `valid:"required"`
	Age       uint8  `valid:"between=0|30"`
	Email     string `valid:"required,email"`
}

type appError struct {
	Code    int              `json:"code,omitempty"`
	Errors  validator.Errors `json:"errors,omitempty"`
	Message string           `json:"message,omitempty"`
}

// ErrorResponse return error
func ErrorResponse(c *gin.Context, code int, err error) {
	switch v := err.(type) {
	case validator.Errors:
		locale := c.DefaultQuery("locale", "en")
		c.JSON(http.StatusBadRequest, gin.H{"error": &appError{
			Code:   http.StatusBadRequest,
			Errors: binding.Validator.Engine().(*validator.Validator).Translator.Trans(v, locale),
		}})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": &appError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}})
	}
	return
}

func init() {

	// replace gin default validator
	binding.Validator = new(DefaultValidator)
	if v, ok := binding.Validator.Engine().(*validator.Validator); ok {
		v.Translator = validator.NewTranslator()
		v.Translator.SetMessage("en", validator_en.MessageMap)
		v.Translator.SetMessage("zh_CN", validator_zhCN.MessageMap)
		v.Translator.SetAttributes("en", lang_en.AttributeMap)
		v.Translator.SetAttributes("zh_CN", lang_zhCN.AttributeMap)
	}
}

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		time.Sleep(time.Duration(5) * time.Second)
		var form User
		if err := c.ShouldBind(&form); err == nil {
			c.JSON(http.StatusOK, &form)
		} else {
			ErrorResponse(c, http.StatusBadRequest, err)
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
