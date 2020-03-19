//+build faker

package monitorable

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/http"
	httpDelivery "github.com/monitoror/monitoror/monitorables/http/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/models"
	httpUsecase "github.com/monitoror/monitoror/monitorables/http/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type httpMonitorable struct{}

func NewHTTPMonitorable(_ map[string]*config.HTTP, _ cache.Store, _ int) Monitorable {
	return &httpMonitorable{}
}

func (m *httpMonitorable) GetHelp() string       { return "" }
func (m *httpMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *httpMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := httpUsecase.NewHTTPUsecase()
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// EnableTile route to echo
	httpGroup := router.Group("/http", variant)
	routeStatus := httpGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := httpGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// EnableTile data for config hydration
	configManager.EnableTile(http.HTTPStatusTileType, variant, &httpModels.HTTPStatusParams{}, routeStatus.Path, config.DefaultInitialMaxDelay)
	configManager.EnableTile(http.HTTPRawTileType, variant, &httpModels.HTTPRawParams{}, routeRaw.Path, config.DefaultInitialMaxDelay)
	configManager.EnableTile(http.HTTPFormattedTileType, variant, &httpModels.HTTPFormattedParams{}, routeJSON.Path, config.DefaultInitialMaxDelay)

	return true
}
