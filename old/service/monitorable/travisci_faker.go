//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/travisci"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/models"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type travisciMonitorable struct{}

func NewTravisCIMonitorable(_ map[string]*config.TravisCI) Monitorable {
	return &travisciMonitorable{}
}

func (m *travisciMonitorable) GetHelp() string       { return "" }
func (m *travisciMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *travisciMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := travisciUsecase.NewTravisCIUsecase()
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	route := router.Group("/travisci", variant).GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	configManager.EnableTile(travisci.TravisCIBuildTileType, variant, &travisciModels.BuildParams{}, route.Path, config.DefaultInitialMaxDelay)

	return true
}
