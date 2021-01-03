package gitlab

import (
	"fmt"

	config "github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	generic "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"

	"strconv"

	log "github.com/sirupsen/logrus"
	gl "github.com/xanzy/go-gitlab"
)

type Project struct {
	ref *gl.Project
}

func (p *Project) ID() string {
	return strconv.Itoa(p.ref.ID)
}
func (p *Project) Name() string {
	return p.ref.Name
}
func (p *Project) IsArchived() bool {
	return p.ref.Archived
}
func (p *Project) BrowserURL() string {
	return p.ref.WebURL
}
func (p *Project) Path() string {
	return p.ref.Namespace.FullPath
}

func (p *Project) GetCloneURI(protocol config.GitAccessProtocol) string {
	switch {
	case protocol == config.SSH:
		return p.ref.SSHURLToRepo
	case protocol == config.HTTP:
		return p.ref.HTTPURLToRepo
	default:
		return p.ref.SSHURLToRepo
	}
}

type fetchProject struct {
	client *gl.Client
}

func (p fetchProject) Get(id string) ([]generic.RemoteRepoInfo, error) {
	log.Debugf("Load Project '%s' from Gitlab", id)

	if projectID, err := strconv.Atoi(id); err == nil {
		project, _, err := p.client.Projects.GetProject(projectID, nil, nil)
		if err != nil {
			return nil, err
		}

		var elements []generic.RemoteRepoInfo

		pr := Project{
			ref: project,
		}
		elements = append(elements, &pr)

		return elements, nil
	}

	return nil, fmt.Errorf("project with %s not found", id)
}
