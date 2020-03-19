//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/azuredevops"
	azuredevopsDelivery "github.com/monitoror/monitoror/monitorables/azuredevops/delivery/http"
	azureDevOpsModels "github.com/monitoror/monitoror/monitorables/azuredevops/models"
	azuredevopsUsecase "github.com/monitoror/monitoror/monitorables/azuredevops/usecase"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/service/router"
)

type azuredevopsMonitorable struct{}

func NewAzureDevOpsMonitorable(_ map[string]*config.AzureDevOps) Monitorable {
	return &azuredevopsMonitorable{}
}

func (m *azuredevopsMonitorable) GetHelp() string       { return "" }
func (m *azuredevopsMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *azuredevopsMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := azuredevopsUsecase.NewAzureDevOpsUsecase()
	delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

	// EnableTile route to echo
	azureGroup := router.Group("/azuredevops", variant)
	routeBuild := azureGroup.GET("/build", delivery.GetBuild)
	routeRelease := azureGroup.GET("/release", delivery.GetRelease)

	// EnableTile data for config hydration
	configManager.EnableTile(azuredevops.AzureDevOpsBuildTileType, variant, &azureDevOpsModels.BuildParams{}, routeBuild.Path, config.DefaultInitialMaxDelay)
	configManager.EnableTile(azuredevops.AzureDevOpsReleaseTileType, variant, &azureDevOpsModels.ReleaseParams{}, routeRelease.Path, config.DefaultInitialMaxDelay)

	return true
}
