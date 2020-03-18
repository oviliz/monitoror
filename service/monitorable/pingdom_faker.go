//+build faker

package monitorable

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/pingdom"
	pingdomDelivery "github.com/monitoror/monitoror/monitorables/pingdom/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/models"
	pingdomUsecase "github.com/monitoror/monitoror/monitorables/pingdom/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type pingdomMonitorable struct{}

func NewPingdomMonitorable(_ map[string]*config.Pingdom, _ cache.Store) Monitorable {
	return &pingdomMonitorable{}
}

func (m *pingdomMonitorable) GetHelp() string       { return "" }
func (m *pingdomMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *pingdomMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := pingdomUsecase.NewPingdomUsecase()
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	route := router.Group("/pingdom", variant).GET("/check", delivery.GetCheck)

	// EnableTile data for config hydration
	configManager.EnableTile(pingdom.PingdomCheckTileType, variant, &pingdomModels.CheckParams{}, route.Path, config.DefaultInitialMaxDelay)
	configManager.DisableTile(pingdom.PingdomChecksTileType, variant)

	return true
}
