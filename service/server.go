package service

import (
	"fmt"
	"math/rand"
	"time"

	"honnef.co/go/tools/config"

	uiConfig "github.com/monitoror/monitoror/api/config"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/handlers"
	"github.com/monitoror/monitoror/middlewares"
	"github.com/monitoror/monitoror/service/router"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Server struct {
		// Echo Server
		*echo.Echo

		store *Store
	}

	// Store is used to share Data in every monitorable
	Store struct {
		// Global CoreConfig
		CoreConfig *coreConfig.Config

		// CacheStore for every memory persistent data
		CacheStore cache.Store
		// MidCacheMiddlewaredleware using CacheStore to return cached data
		CacheMiddleware *middlewares.CacheMiddleware

		// MonitorableRouter helper wrapping echo Router monitorable
		MonitorableRouter router.MonitorableRouter

		// MonitorableConfigManager used to register Tile for verify / hydrate
		UiConfigManager uiConfig.Manager
	}
)

var colorer = color.New()

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Init create echo server with middlewares, ui, routes
func Init(config *config.Config) *Server {
	server := &Server{
		store: &Store{
			CoreConfig: config,
		},
	}

	server.setupEchoServer()
	server.setupEchoMiddleware()

	InitUI(server)
	InitApis(server)

	return server
}

func (s *Server) Start() {
	fmt.Println()
	log.Fatal(s.Echo.Start(fmt.Sprintf(":%d", s.store.CoreConfig.Port)))
}

func (s *Server) setupEchoServer() {
	s.Echo = echo.New()
	s.HideBanner = true

	// ----- Errors Handler -----
	s.HTTPErrorHandler = handlers.HTTPErrorHandler
}

func (s *Server) setupEchoMiddleware() {
	// Recover (don't panic 😎)
	s.Use(echoMiddleware.Recover())

	// Log requests
	if s.store.CoreConfig.Env != "production" {
		s.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format: `[-] ` + colorer.Green("${method}") + ` ${uri} status:${status} latency:` + colorer.Green("${latency_human}") + ` error:"${error}"` + "\n",
		}))
	}

	// Cache
	s.store.CacheStore = cache.NewGoCacheStore(time.Minute*5, time.Second) // Default value, always override
	s.store.CacheMiddleware = middlewares.NewCacheMiddleware(s.store.CacheStore,
		time.Millisecond*time.Duration(s.store.CoreConfig.DownstreamCacheExpiration),
		time.Millisecond*time.Duration(s.store.CoreConfig.UpstreamCacheExpiration),
	) // Used as Handler wrapper in routes
	s.Use(s.store.CacheMiddleware.DownstreamStoreMiddleware())

	// CORS
	s.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
}
