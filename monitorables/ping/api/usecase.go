package api

import (
	"github.com/monitoror/monitoror/models"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/api/models"
)

const (
	PingTileType models.TileType = "PING"
)

type (
	Usecase interface {
		Ping(params *pingModels.PingParams) (*models.Tile, error)
	}
)
