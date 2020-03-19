package http

import (
	"github.com/monitoror/monitoror/monitorables/http/models"
)

type (
	Repository interface {
		Get(url string) (*models.Response, error)
	}
)
