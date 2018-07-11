package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/syssam/go-validator"
	lang_en "github.com/syssam/go-validator/examples/translations/lang/en"
	lang_zhCN "github.com/syssam/go-validator/examples/translations/lang/zh_CN"
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

func init() {

	// replace gin default validator
	binding.Validator = new(DefaultValidator)
	if v, ok := binding.Validator.Engine().(*validator.Validator); ok {
		v.Translator = validator.NewTranslator()
		v.Translator.SetMessage("en", validator_en.MessageMap)
		v.Translator.SetMessage("zh_CN", validator_zhCN.MessageMap)
		v.Translator.SetAttributes("en", lang_en.AttributeMap)
		v.Translator.SetAttributes("zh_CN", lang_zhCN.AttributeMap)
		v.Translator.SetLocale("zh_CN")
	}

}

// LocalizationMiddleware init
func LocalizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.DefaultQuery("locale", "en")
		if v, ok := binding.Validator.Engine().(*validator.Validator); ok {
			v.Translator.SetLocale(locale)
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(LocalizationMiddleware())
	r.POST("/", func(c *gin.Context) {
		var form User
		if err := c.ShouldBind(&form); err == nil {
			c.JSON(http.StatusOK, &form)
		} else {
			if err.(validator.Errors) != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.(validator.Errors)})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
