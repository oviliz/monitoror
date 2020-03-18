package pingdom

import (
	"github.com/monitoror/monitoror/monitorables/pingdom/models"
)

type (
	Repository interface {
		GetCheck(checkID int) (*models.Check, error)
		GetChecks(tags string) ([]models.Check, error)
	}
)
