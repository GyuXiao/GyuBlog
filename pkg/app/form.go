// Package app
/**
  @author: zk.xiao
  @date: 2022/5/18
  @note:
**/
package app

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"strings"
)

type ValidError struct {
	key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ", ")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors

	// 参数绑定 + 入参校验
	err := c.ShouldBind(v)

	// 发生错误后，通过在中间件 Translations 设置的 Translator 来对错误消息体进行具体的翻译行为
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}
		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				key:     key,
				Message: value,
			})
		}
		return false, errs
	}
	return true, nil
}
