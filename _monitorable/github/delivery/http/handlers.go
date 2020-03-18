package http

import (
	"net/http"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github"
	githubModels "github.com/monitoror/monitoror/monitorables/github/models"

	"github.com/labstack/echo/v4"
)

type GithubDelivery struct {
	githubUsecase github.Usecase
}

func NewGithubDelivery(p github.Usecase) *GithubDelivery {
	return &GithubDelivery{p}
}

func (h *GithubDelivery) GetCount(c echo.Context) error {
	// Bind / check Params
	params := &githubModels.CountParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.githubUsecase.Count(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}

func (h *GithubDelivery) GetChecks(c echo.Context) error {
	// Bind / check Params
	params := &githubModels.ChecksParams{}
	err := c.Bind(params)
	if err != nil || !params.IsValid() {
		return models.QueryParamsError
	}

	tile, err := h.githubUsecase.Checks(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tile)
}
