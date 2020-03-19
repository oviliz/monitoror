package jenkins

import (
	"github.com/monitoror/monitoror/models"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	JenkinsBuildTileType       models.TileType = "JENKINS-BUILD"
	JenkinsMultiBranchTileType models.TileType = "JENKINS-MULTIBRANCH"
)

type (
	Usecase interface {
		Build(params *jenkinsModels.BuildParams) (*models.Tile, error)
		MultiBranch(params interface{}) ([]builder.Result, error)
	}
)
