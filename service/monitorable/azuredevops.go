//+build !faker

package monitorable

import (
	"net/url"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/azuredevops"
	azuredevopsDelivery "github.com/monitoror/monitoror/monitorables/azuredevops/delivery/http"
	azureDevOpsModels "github.com/monitoror/monitoror/monitorables/azuredevops/models"
	azuredevopsRepository "github.com/monitoror/monitoror/monitorables/azuredevops/repository"
	azuredevopsUsecase "github.com/monitoror/monitoror/monitorables/azuredevops/usecase"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/service/router"
)

type azuredevopsMonitorable struct {
	config map[string]*config.AzureDevOps
}

func NewAzureDevOpsMonitorable(config map[string]*config.AzureDevOps) Monitorable {
	return &azuredevopsMonitorable{config: config}
}

func (m *azuredevopsMonitorable) GetHelp() string { return "HEEEELLLPPPP" }
func (m *azuredevopsMonitorable) GetVariants() []string {
	return config.GetVariantsFromConfig(m.config)
}
func (m *azuredevopsMonitorable) isEnabled(variant string) bool {
	conf := m.config[variant]

	if conf.URL == "" {
		return false
	}

	if _, err := url.Parse(conf.URL); err != nil {
		return false
	}

	return conf.Token != ""
}

func (m *azuredevopsMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]
		repository := azuredevopsRepository.NewAzureDevOpsRepository(conf)
		usecase := azuredevopsUsecase.NewAzureDevOpsUsecase(repository)
		delivery := azuredevopsDelivery.NewAzureDevOpsDelivery(usecase)

		// EnableTile route to echo
		azureGroup := router.Group("/azuredevops", variant)
		routeBuild := azureGroup.GET("/build", delivery.GetBuild)
		routeRelease := azureGroup.GET("/release", delivery.GetRelease)

		// EnableTile data for config hydration
		configManager.EnableTile(azuredevops.AzureDevOpsBuildTileType, variant, &azureDevOpsModels.BuildParams{}, routeBuild.Path, conf.InitialMaxDelay)
		configManager.EnableTile(azuredevops.AzureDevOpsReleaseTileType, variant, &azureDevOpsModels.ReleaseParams{}, routeRelease.Path, conf.InitialMaxDelay)
	} else {
		// EnableTile data for config verify
		configManager.DisableTile(azuredevops.AzureDevOpsBuildTileType, variant)
		configManager.DisableTile(azuredevops.AzureDevOpsReleaseTileType, variant)
	}

	return enabled
}
