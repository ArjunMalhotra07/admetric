package http

import (
	"encoding/json"

	"github.com/ArjunMalhotra/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App
	Log *logger.Logger
}

func NewApp(log *logger.Logger) *App {
	newapp := fiber.New(fiber.Config{
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		EnableTrustedProxyCheck: true,
	})

	return &App{
		App: newapp,
		Log: log,
	}
}

type HttpResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
}

const (
	StatusBadRequest          = fiber.StatusBadRequest
	StatusUnauthorized        = fiber.StatusUnauthorized
	StatusForbidden           = fiber.StatusForbidden
	StatusNotFound            = fiber.StatusNotFound
	StatusInternalServerError = fiber.StatusInternalServerError
	StatusOK                  = fiber.StatusOK
	StatusCreated             = fiber.StatusCreated
	StatusNoContent           = fiber.StatusNoContent
)

const (
	ErrBadRequest          = "Bad request"
	ErrInternalServerError = "Internal server error"
	ErrAlreadyExists       = "Already exists"
	ErrNotFound            = "Not Found"
	ErrUnauthorized        = "Unauthorized"
	ErrForbidden           = "Forbidden"
	ErrBadQueryParams      = "Invalid query params"
	ErrRequestTimeout      = "Request Timeout"
	ErrEndpointNotFound    = "The endpoint you requested doesn't exist on server"
)

// http 200 ok http response
func (a *App) HttpResponseOK(c *fiber.Ctx, data interface{}) error {
	return c.Status(StatusOK).JSON(
		&HttpResponse{
			Success: true,
			Code:    StatusOK,
			Data:    data,
			Error:   "",
			Message: "",
		})
}

// http 201 created http response
func (a *App) HttpResponseCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(StatusCreated).JSON(
		&HttpResponse{
			Success: true,
			Code:    StatusCreated,
			Data:    data,
			Error:   "",
			Message: "",
		})
}

// http 204 no content http response
func (a *App) HttpResponseNoContent(c *fiber.Ctx) error {
	return c.Status(StatusNoContent).JSON(
		&HttpResponse{
			Success: true,
			Code:    StatusNoContent,
			Data:    nil,
			Error:   "",
			Message: "",
		})
}

// http 400 bad request http response
func (a *App) HttpResponseBadRequest(c *fiber.Ctx, message error) error {
	a.Log.Logger.Error(message.Error())
	return c.Status(StatusBadRequest).JSON(
		&HttpResponse{
			Success: false,
			Code:    StatusBadRequest,
			Data:    nil,
			Error:   ErrBadRequest,
			Message: message.Error(),
		})
}

// http 400 bad query params http response
func (a *App) HttpResponseBadQueryParams(c *fiber.Ctx, message error) error {
	a.Log.Logger.Error(message.Error())
	return c.Status(StatusBadRequest).JSON(
		&HttpResponse{
			Success: false,
			Code:    StatusBadRequest,
			Data:    nil,
			Error:   ErrBadQueryParams,
			Message: message.Error(),
		})
}

// http 404 not found http response
func (a *App) HttpResponseNotFound(c *fiber.Ctx, message error) error {
	a.Log.Logger.Error(message.Error())
	return c.Status(StatusNotFound).JSON(
		&HttpResponse{
			Success: false,
			Code:    StatusNotFound,
			Data:    nil,
			Error:   ErrNotFound,
			Message: message.Error(),
		})
}

// http 500 internal server error response
func (a *App) HttpResponseInternalServerErrorRequest(c *fiber.Ctx, message error) error {
	a.Log.Logger.Error(message.Error())
	return c.Status(StatusInternalServerError).JSON(
		&HttpResponse{
			Success: false,
			Code:    StatusInternalServerError,
			Data:    nil,
			Error:   ErrInternalServerError,
			Message: message.Error(),
		})
}

// http 403 The client does not have access rights to the content;
// that is, it is unauthorized, so the server is refusing to give the requested resource
func (a *App) HttpResponseForbidden(c *fiber.Ctx, message error) error {
	a.Log.Logger.Error(message.Error())
	return c.Status(StatusForbidden).JSON(
		&HttpResponse{
			Success: false,
			Code:    StatusForbidden,
			Data:    nil,
			Error:   ErrForbidden,
			Message: message.Error(),
		})
}

// http 401 the client must authenticate itself to get the requested response
func (a *App) HttpResponseUnauthorized(c *fiber.Ctx, message error) error {
	a.Log.Logger.Error(message.Error())
	return c.Status(StatusUnauthorized).JSON(
		&HttpResponse{
			Success: false,
			Code:    StatusUnauthorized,
			Data:    nil,
			Error:   ErrUnauthorized,
			Message: message.Error(),
		})
}
