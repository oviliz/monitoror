package github

import (
	"github.com/monitoror/monitoror/models"
	githubModels "github.com/monitoror/monitoror/monitorables/github/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	GithubCountTileType       models.TileType = "GITHUB-COUNT"
	GithubChecksTileType      models.TileType = "GITHUB-CHECKS"
	GithubPullRequestTileType models.TileType = "GITHUB-PULLREQUESTS"
)

type (
	Usecase interface {
		Count(params *githubModels.CountParams) (*models.Tile, error)
		Checks(params *githubModels.ChecksParams) (*models.Tile, error)
		PullRequests(params interface{}) ([]builder.Result, error)
	}
)
