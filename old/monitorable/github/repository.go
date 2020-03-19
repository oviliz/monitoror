package github

import "github.com/monitoror/monitoror/monitorables/github/models"

type (
	Repository interface {
		GetCount(query string) (int, error)
		GetChecks(owner, repository, ref string) (*models.Checks, error)
		GetPullRequests(owner, repository string) ([]models.PullRequest, error)
		GetCommit(owner, repository, sha string) (*models.Commit, error)
	}
)
