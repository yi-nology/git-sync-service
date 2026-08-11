package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
}

func Error(c *app.RequestContext, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Code:    statusCode,
		Message: message,
	})
}

func BadRequest(c *app.RequestContext, message string) {
	Error(c, consts.StatusBadRequest, message)
}

func InternalError(c *app.RequestContext, message string) {
	Error(c, consts.StatusInternalServerError, message)
}

func NotFound(c *app.RequestContext, message string) {
	Error(c, consts.StatusNotFound, message)
}

func Success(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    consts.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Created(c *app.RequestContext, data any) {
	c.JSON(consts.StatusCreated, SuccessResponse{
		Code:    consts.StatusCreated,
		Message: "created",
		Data:    data,
	})
}

func NoContent(c *app.RequestContext) {
	c.Status(consts.StatusNoContent)
}
