package travisci

import (
	"github.com/monitoror/monitoror/monitorables/travisci/models"
)

type (
	Repository interface {
		GetLastBuildStatus(owner, repository, branch string) (*models.Build, error)
	}
)
