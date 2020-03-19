package config

import (
	"fmt"
	"reflect"
	"strings"

	pkgConfig "github.com/monitoror/monitoror/pkg/monitoror/config"

	"github.com/spf13/viper"
)

const EnvPrefix = "MO"
const MonitorablePrefix = "MONITORABLE"

const DefaultVariant = "default"
const DefaultInitialMaxDelay = 1700

type (
	// CoreConfig contain backend Configuration
	Config struct {
		// --- General Configuration ---
		Port int    // Default: 8080
		Env  string // Default: production

		// --- Cache Configuration ---
		// UpstreamCacheExpiration is used to respond before executing the request. Avoid overloading services.
		UpstreamCacheExpiration int
		// DownstreamCacheExpiration is used to respond after executing the request in case of timeout error.
		DownstreamCacheExpiration int

		// Monitorable CoreConfig
		Monitorable Monitorable
	}

	Monitorable struct {
		Ping        map[string]*Ping
		Port        map[string]*Port
		HTTP        map[string]*HTTP
		Pingdom     map[string]*Pingdom
		TravisCI    map[string]*TravisCI
		Jenkins     map[string]*Jenkins
		AzureDevOps map[string]*AzureDevOps
		Github      map[string]*Github
	}

	Port struct {
		Timeout         int // In Millisecond
		InitialMaxDelay int // In Millisecond
	}

	HTTP struct {
		Timeout         int // In Millisecond
		SSLVerify       bool
		InitialMaxDelay int // In Millisecond
	}

	Pingdom struct {
		URL             string
		Token           string
		Timeout         int // In Millisecond
		CacheExpiration int // In Millisecond
		InitialMaxDelay int // In Millisecond
	}

	TravisCI struct {
		URL             string
		Timeout         int // In Millisecond
		Token           string
		GithubToken     string
		InitialMaxDelay int // In Millisecond
	}

	Jenkins struct {
		URL             string
		Timeout         int // In Millisecond
		SSLVerify       bool
		Login           string
		Token           string
		InitialMaxDelay int // In Millisecond
	}

	AzureDevOps struct {
		URL             string
		Timeout         int // In Millisecond
		Token           string
		InitialMaxDelay int // In Millisecond
	}

	Github struct {
		Timeout              int // In Millisecond
		Token                string
		CountCacheExpiration int // In Millisecond
		InitialMaxDelay      int // In Millisecond
	}
)

// InitConfig from configuration file / env / default value
func InitConfig() *Config {
	var config Config

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Transform Env and define Label for setting default value
	variants := initEnvAndVariant()

	// Setup default values
	// --- General Configuration ---
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Env", "production")

	// --- Cache Configuration ---
	viper.SetDefault("UpstreamCacheExpiration", 10000)
	viper.SetDefault("DownstreamCacheExpiration", 120000)

	// --- Ping Configuration ---
	for variant := range variants["Ping"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.Count", variant), 2)
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.Timeout", variant), 1000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.Interval", variant), 100)
		viper.SetDefault(fmt.Sprintf("Monitorable.Ping.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- Port Configuration ---
	for variant := range variants["Port"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Port.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Port.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- HTTP Configuration ---
	for variant := range variants["HTTP"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.SSLVerify", variant), true)
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.URL", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.HTTP.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- Pingdom Configuration ---
	for variant := range variants["Pingdom"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.URL", variant), "https://api.pingdom.com/api/3.1")
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.CacheExpiration", variant), 30000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Pingdom.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- TravisCI Configuration ---
	for variant := range variants["TravisCI"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.URL", variant), "https://api.travis-ci.com/")
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.TravisCI.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- Jenkins Configuration ---
	for variant := range variants["Jenkins"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.URL", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Timeout", variant), 2000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.SSLVerify", variant), true)
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Login", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Jenkins.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- Azure DevOps Configuration ---
	for variant := range variants["AzureDevOps"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.URL", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.Timeout", variant), 4000)
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.AzureDevOps.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	// --- Github Configuration ---
	for variant := range variants["Github"] {
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.Timeout", variant), 5000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.Token", variant), "")
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.CountCacheExpiration", variant), 30000)
		viper.SetDefault(fmt.Sprintf("Monitorable.Github.%s.InitialMaxDelay", variant), DefaultInitialMaxDelay)
	}

	_ = viper.Unmarshal(&config)

	return &config
}

// -------- Config Utility function ---------
func LoadMonitorableConfig(conf interface{}, defaultConf interface{}) {
	pkgConfig.LoadConfig(fmt.Sprintf("%s_%s", EnvPrefix, MonitorablePrefix), DefaultVariant, conf, defaultConf)
}

func GetVariantsFromConfig(conf interface{}) []string {
	var variants []string
	if reflect.TypeOf(conf).Kind() == reflect.Map {
		keys := reflect.ValueOf(conf).MapKeys()
		for _, k := range keys {
			variants = append(variants, k.String())
		}
	} else {
		variants = append(variants, DefaultVariant)
	}

	return variants
}
