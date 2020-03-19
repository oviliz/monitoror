package http

import (
	netHttp "net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/models"

	"github.com/labstack/echo/v4"
)

//nolint:golint
type HTTPDelivery struct {
	httpUsecase http.Usecase
}

func NewHTTPDelivery(p http.Usecase) *HTTPDelivery {
	return &HTTPDelivery{p}
}

func (h *HTTPDelivery) GetHTTPStatus(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPStatusParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPStatus(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPRaw(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPRawParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPRaw(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}

func (h *HTTPDelivery) GetHTTPFormatted(c echo.Context) error {
	// Bind / Check Params
	params := &httpModels.HTTPFormattedParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.httpUsecase.HTTPFormatted(params)
	if err != nil {
		return err
	}

	return c.JSON(netHttp.StatusOK, tile)
}
