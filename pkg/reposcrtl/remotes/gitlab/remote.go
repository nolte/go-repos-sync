package gitlab

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	generic "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	gl "github.com/xanzy/go-gitlab"
)

func GetServices(cfg config.RemoteEndpoint) (map[config.DirectoryType]generic.RemoteDataLookup, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, err
	}

	services := map[config.DirectoryType]generic.RemoteDataLookup{
		config.DirectoryTypeGroup:   fetchGroupProjects{client: client},
		config.DirectoryTypeProject: fetchProject{client: client},
	}

	return services, nil
}

func newClient(cfg config.RemoteEndpoint) (*gl.Client, error) {
	tokenValue, err := cfg.Auth.GetTokenValue()
	if err != nil {
		return nil, err
	}

	if cfg.URI != "" {
		client, err := gl.NewClient(tokenValue, gl.WithBaseURL(cfg.URI))
		return client, err
	}

	return gl.NewClient(tokenValue)
}
