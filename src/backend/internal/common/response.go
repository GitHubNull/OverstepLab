package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess       = 0
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeInternalError = 500
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func SuccessMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: msg,
	})
}

func Error(c *gin.Context, httpCode, code int, msg string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: msg,
	})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, CodeForbidden, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, CodeNotFound, msg)
}

func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, msg)
}
