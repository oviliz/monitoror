//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/github"
	githubDelivery "github.com/monitoror/monitoror/monitorables/github/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorables/github/models"
	githubUsecase "github.com/monitoror/monitoror/monitorables/github/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type githubMonitorable struct{}

func NewGithubMonitorable(_ map[string]*config.Github) Monitorable { return &githubMonitorable{} }

func (m *githubMonitorable) GetHelp() string       { return "" }
func (m *githubMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *githubMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := githubUsecase.NewGithubUsecase()
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// EnableTile route to echo
	azureGroup := router.Group("/github", variant)
	routeCount := azureGroup.GET("/count", delivery.GetCount)
	routeChecks := azureGroup.GET("/checks", delivery.GetChecks)

	// EnableTile data for config hydration
	configManager.EnableTile(github.GithubCountTileType, variant, &githubModels.CountParams{}, routeCount.Path, config.DefaultInitialMaxDelay)
	configManager.EnableTile(github.GithubChecksTileType, variant, &githubModels.ChecksParams{}, routeChecks.Path, config.DefaultInitialMaxDelay)
	configManager.DisableTile(github.GithubPullRequestTileType, variant)

	return true
}
