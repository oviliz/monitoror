package ping

import (
	"github.com/monitoror/monitoror/monitorables/ping/models"
)

type (
	Repository interface {
		ExecutePing(hostname string) (*models.Ping, error)
	}
)
