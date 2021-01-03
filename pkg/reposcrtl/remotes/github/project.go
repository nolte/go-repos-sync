package github

import (
	"context"
	"strings"

	gh "github.com/google/go-github/github"
	config "github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	generic "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
)

type Project struct {
	ref *gh.Repository
}

func (p *Project) ID() string {
	return *p.ref.Name
}
func (p *Project) Name() string {
	return *p.ref.Name
}
func (p *Project) IsArchived() bool {
	return *p.ref.Archived
}
func (p *Project) Path() string {
	return ""
}

func (p *Project) BrowserURL() string {
	return *p.ref.HTMLURL
}
func (p *Project) GetCloneURI(protocol config.GitAccessProtocol) string {
	switch {
	case protocol == config.SSH:
		return *p.ref.SSHURL
	case protocol == config.HTTP:
		return *p.ref.CloneURL
	default:
		return *p.ref.SSHURL
	}
}

type fetchProject struct {
	client *gh.Client
}

func (p fetchProject) Get(id string) ([]generic.RemoteRepoInfo, error) {
	repoParts := strings.Split(id, "/")

	repo, _, errPr := p.client.Repositories.Get(context.Background(), repoParts[0], repoParts[1])
	if errPr != nil {
		return nil, errPr
	}

	var elements []generic.RemoteRepoInfo
	elements = append(elements, &Project{
		ref: repo,
	})

	return elements, errPr
}
