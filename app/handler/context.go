package handler

import (
	"errors"
	"net/http"

	"github.com/CyrivlClth/kube-go/app/errcode"
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) AbortWithError(err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, errcode.ErrBadRequest) {
		code = http.StatusBadRequest
	}
	c.Context.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
}

func (c *Context) JSON(data any) {
	c.Context.JSON(http.StatusOK, gin.H{"data": data})
}

func WrapHandle(fn func(*Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{Context: c}
		err := fn(ctx)
		if err != nil {
			ctx.AbortWithError(err)
		}
	}
}
