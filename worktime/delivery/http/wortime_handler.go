package http

import (
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/labstack/echo"
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
	route := e.Group("ayaya")
	route.GET("/worktimes", handler.FetchWorktime)
}

func (handler *WorktimeHandler) FetchWorktime(ec echo.Context) error {
	ctx := ec.Request().Context()
	result, err := handler.WorktimeUsecase.Fetch(ctx)

	if err != nil {
		pp.Println(err)
		return ec.JSON(http.StatusInternalServerError, err)
	}

	response := make(map[string]interface{})
	response["data"] = result

	return ec.JSON(http.StatusOK, response)
}
