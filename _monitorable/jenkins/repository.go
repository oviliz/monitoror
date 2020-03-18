package jenkins

import (
	"github.com/monitoror/monitoror/monitorables/jenkins/models"
)

type (
	Repository interface {
		GetJob(jobName string, branch string) (*models.Job, error)
		GetLastBuildStatus(job *models.Job) (*models.Build, error)
	}
)
