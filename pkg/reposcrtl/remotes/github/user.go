package github

import (
	"context"

	gh "github.com/google/go-github/github"
	generic "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
)

type fetchUserProjects struct {
	client *gh.Client
}

const ghMaxPerPage = 30

func (p fetchUserProjects) Get(id string) ([]generic.RemoteRepoInfo, error) {
	opt := &gh.RepositoryListOptions{
		ListOptions: gh.ListOptions{PerPage: ghMaxPerPage},
	}

	// get all pages of results
	var allRepos []generic.RemoteRepoInfo

	for {
		repos, resp, err := p.client.Repositories.List(context.Background(), id, opt)
		if err != nil {
			return nil, err
		}

		for _, element := range repos {
			allRepos = append(allRepos, &Project{ref: element})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allRepos, nil
}
