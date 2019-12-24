package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	basicAuth "github.com/laminh711/reporter/auth"

	"github.com/labstack/echo/middleware"
	"github.com/laminh711/reporter/models"
	"github.com/laminh711/reporter/worklog"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Data interface{} `json:"data"`
}

type CreateRequest struct {
	Name       string `json:"name"`
	Productive bool   `json:"productive"`
}

type WorklogHandler struct {
	WorklogUsecase worklog.Usecase
}

func NewWorklogHandler(e *echo.Echo, wu worklog.Usecase) {
	handler := &WorklogHandler{
		WorklogUsecase: wu,
	}

	route := e.Group("api")
	route.Use(middleware.BasicAuth(basicAuth.Middleware))

	route.GET("/worklogs", handler.FetchWorklog)
	route.POST("/worklogs", handler.CreateWorklog)
}

func (handler *WorklogHandler) FetchWorklog(ec echo.Context) error {
	ctx := ec.Request().Context()
	result, err := handler.WorklogUsecase.Fetch(ctx)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
	}
	return ec.JSON(http.StatusOK, ResponseSuccess{result})
}

func (handler *WorklogHandler) CreateWorklog(ec echo.Context) error {
	request := new(CreateRequest)
	if err := ec.Bind(request); err != nil {
		return ec.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
	}

	worklog := models.Worklog{
		Name:       request.Name,
		Productive: request.Productive,
		FinishedAt: time.Now().UTC(),
	}
	ctx := ec.Request().Context()
	err := handler.WorklogUsecase.Create(ctx, worklog)
	if err != nil {
		return ec.JSON(http.StatusInternalServerError, ResponseError{err.Error()})
	}

	return ec.JSON(http.StatusOK, ResponseSuccess{"ok"})
}
