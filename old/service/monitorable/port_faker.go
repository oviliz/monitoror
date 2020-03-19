//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/port"
	portDelivery "github.com/monitoror/monitoror/monitorables/port/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorables/port/models"
	portUsecase "github.com/monitoror/monitoror/monitorables/port/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type portMonitorable struct{}

func NewPortMonitorable(_ map[string]*config.Port) Monitorable {
	return &portMonitorable{}
}

func (m *portMonitorable) GetHelp() string       { return "" }
func (m *portMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *portMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := portUsecase.NewPortUsecase()
	delivery := portDelivery.NewPortDelivery(usecase)

	// EnableTile route to echo
	route := router.Group("/port", variant).GET("/port", delivery.GetPort)

	// EnableTile data for config hydration
	configManager.EnableTile(port.PortTileType, variant, &portModels.PortParams{}, route.Path, config.DefaultInitialMaxDelay)

	return true
}
