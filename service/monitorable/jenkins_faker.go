//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/jenkins"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorables/jenkins/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/models"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorables/jenkins/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type jenkinsMonitorable struct{}

func NewJenkinsMonitorable(_ map[string]*config.Jenkins) Monitorable {
	return &jenkinsMonitorable{}
}

func (m *jenkinsMonitorable) GetHelp() string       { return "" }
func (m *jenkinsMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *jenkinsMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := jenkinsUsecase.NewJenkinsUsecase()
	delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

	// EnableTile route to echo
	route := router.Group("/http", variant).GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	configManager.EnableTile(jenkins.JenkinsBuildTileType, variant, &jenkinsModels.BuildParams{}, route.Path, config.DefaultInitialMaxDelay)
	configManager.DisableTile(jenkins.JenkinsMultiBranchTileType, variant)

	return true
}
