package e

import (
	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, httpCode int, code ERRCode, field string) {
	WriteResponse(c, httpCode, code, field, nil)
}

func WriteResponse(c *gin.Context, httpCode int, code ERRCode, field string, data interface{}) {
	resp := NewResponse(code)
	if len(field) > 0 {
		resp.Field = field
	}
	if data != nil {
		resp.Data = data
	} else {
		resp.Data = struct{}{}
	}
	c.JSON(httpCode, resp)
}

func WriteListError(c *gin.Context, httpCode int, code ERRCode, field string) {
	WriteListResponse(c, httpCode, code, field, 0, 0, 0, nil)
}

func WriteListResponse(c *gin.Context, httpCode int, code ERRCode, field string, offset, limit, total int64, data interface{}) {
	resp := NewListResponse(code, offset, limit, total)
	if len(field) > 0 {
		resp.Field = field
	}
	if data != nil {
		resp.Data = data
	} else {
		resp.Data = []struct{}{}
	}
	c.JSON(httpCode, resp)
}
