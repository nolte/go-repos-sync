package github

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	generic "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"

	"context"

	gh "github.com/google/go-github/github"

	"net/http"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
	"golang.org/x/oauth2"
)

func GetServices(cfg config.RemoteEndpoint) (map[config.DirectoryType]generic.RemoteDataLookup, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, err
	}

	services := map[config.DirectoryType]generic.RemoteDataLookup{
		config.DirectoryTypeProject: fetchProject{client: client},
		config.DirectoryTypeUser:    fetchUserProjects{client: client},
	}

	return services, nil
}
func newClient(cfg config.RemoteEndpoint) (*gh.Client, error) {
	var tc *http.Client

	if cfg.Auth.Token != "" {
		tokenValue, err := utils.LookupValueByRef(cfg.Auth.Token)
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: tokenValue},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	return gh.NewClient(tc), nil
}
