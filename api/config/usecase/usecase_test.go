package usecase

import (
	"io/ioutil"
	"strings"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/api/config/repository"
	. "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/jenkins"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/models"
	"github.com/monitoror/monitoror/monitorables/ping"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/models"
	"github.com/monitoror/monitoror/monitorables/pingdom"
	pindomModels "github.com/monitoror/monitoror/monitorables/pingdom/models"
	"github.com/monitoror/monitoror/monitorables/port"
	portModels "github.com/monitoror/monitoror/monitorables/port/models"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

func initConfigUsecase(repository config.Repository, store cache.Store) *configUsecase {
	usecase := NewConfigUsecase(repository, store, 5000)

	usecase.EnableTile(ping.PingTileType, DefaultVariant, &pingModels.PingParams{}, "/ping", 1000)
	usecase.EnableTile(port.PortTileType, DefaultVariant, &portModels.PortParams{}, "/port", 1000)
	usecase.EnableTile(jenkins.JenkinsBuildTileType, DefaultVariant, &jenkinsModels.BuildParams{}, "/jenkins/default", 1000)
	usecase.EnableTile(pingdom.PingdomCheckTileType, DefaultVariant, &pindomModels.CheckParams{}, "/pingdom/default", 1000)

	return usecase.(*configUsecase)
}

func readConfig(input string) (configBag *models.ConfigBag, err error) {
	configBag = &models.ConfigBag{}
	reader := ioutil.NopCloser(strings.NewReader(input))
	configBag.Config, err = repository.ReadConfig(reader)
	return
}
