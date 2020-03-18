//+build !faker

package monitorable

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorables/config"
	"github.com/monitoror/monitoror/monitorables/http"
	httpDelivery "github.com/monitoror/monitoror/monitorables/http/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/models"
	httpRepository "github.com/monitoror/monitoror/monitorables/http/repository"
	httpUsecase "github.com/monitoror/monitoror/monitorables/http/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type httpMonitorable struct {
	config map[string]*config.HTTP

	// store used for caching request on same url
	store           cache.Store
	cacheExpiration int
}

func NewHTTPMonitorable(config map[string]*config.HTTP, store cache.Store, cacheExpiration int) Monitorable {
	return &httpMonitorable{config: config, store: store, cacheExpiration: cacheExpiration}
}

func (m *httpMonitorable) GetHelp() string         { return "HEEEELLLPPPP" }
func (m *httpMonitorable) GetVariants() []string   { return config.GetVariantsFromConfig(m.config) }
func (m *httpMonitorable) isEnabled(_ string) bool { return true }

func (m *httpMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		repository := httpRepository.NewHTTPRepository(conf)
		usecase := httpUsecase.NewHTTPUsecase(repository, m.store, m.cacheExpiration)
		delivery := httpDelivery.NewHTTPDelivery(usecase)

		// EnableTile route to echo
		httpGroup := router.Group("/http", variant)
		routeStatus := httpGroup.GET("/status", delivery.GetHTTPStatus)
		routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
		routeJSON := httpGroup.GET("/formatted", delivery.GetHTTPFormatted)

		// EnableTile data for config hydration
		configManager.EnableTile(http.HTTPStatusTileType, variant, &httpModels.HTTPStatusParams{}, routeStatus.Path, conf.InitialMaxDelay)
		configManager.EnableTile(http.HTTPRawTileType, variant, &httpModels.HTTPRawParams{}, routeRaw.Path, conf.InitialMaxDelay)
		configManager.EnableTile(http.HTTPFormattedTileType, variant, &httpModels.HTTPFormattedParams{}, routeJSON.Path, conf.InitialMaxDelay)
	} else {
		// EnableTile data for config verify
		configManager.DisableTile(http.HTTPStatusTileType, variant)
		configManager.DisableTile(http.HTTPRawTileType, variant)
		configManager.DisableTile(http.HTTPFormattedTileType, variant)
	}

	return enabled
}
