//+build !faker

package usecase

import (
	"fmt"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/models"
)

type (
	pingUsecase struct {
		repository api.Repository
	}
)

func NewPingUsecase(repository api.Repository) api.Usecase {
	return &pingUsecase{repository}
}

func (pu *pingUsecase) Ping(params *pingModels.PingParams) (tile *models.Tile, err error) {
	tile = models.NewTile(api.PingTileType)
	tile.Label = params.Hostname

	ping, err := pu.repository.ExecutePing(params.Hostname)
	if err == nil {
		tile.Status = models.SuccessStatus
		tile.WithValue(models.MillisecondUnit)
		tile.Value.Values = append(tile.Value.Values, fmt.Sprintf("%d", ping.Average.Milliseconds()))
	} else {
		tile.Status = models.FailedStatus
		err = nil
	}

	return
}
