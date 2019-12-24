package http

import (
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	basicAuth "github.com/laminh711/reporter/auth"
	"github.com/laminh711/reporter/worktime"
)

// ResponseError response error
type ResponseError struct {
	Message string `json:"message"`
}

// WorktimeHandler worktime handler
type WorktimeHandler struct {
	WorktimeUsecase worktime.Usecase
}

func NewWorktimeHandler(e *echo.Echo, u worktime.Usecase) {
	handler := &WorktimeHandler{
		WorktimeUsecase: u,
	}
	route := e.Group("api")
	route.Use(middleware.BasicAuth(basicAuth.Middleware))
	route.GET("/worktimes", handler.FetchWorktime)
}

func (handler *WorktimeHandler) FetchWorktime(ec echo.Context) error {
	ctx := ec.Request().Context()

	pp.Println(ec.Request().Header)

	result, err := handler.WorktimeUsecase.Fetch(ctx)

	if err != nil {
		pp.Println(err)
		return ec.JSON(http.StatusInternalServerError, err)
	}

	response := make(map[string]interface{})
	response["data"] = result

	return ec.JSON(http.StatusOK, response)
}
