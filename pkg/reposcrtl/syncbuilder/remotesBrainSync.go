package syncbuilder

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	github "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes/github"
	gitlab "github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes/gitlab"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"
)

func NewRemotesHostsSyncFunction(db *gorm.DB) *RemotesHostsSyncFunction {
	return &RemotesHostsSyncFunction{
		db: db,
	}
}

type RemotesHostsFactory struct {
	Remotes map[string]*remotes.RemoteFacade
}

func (r *RemotesHostsFactory) GetFacade(alias string) (*remotes.RemoteFacade, error) {
	facade := r.Remotes[alias]
	if facade == nil {
		return nil, fmt.Errorf("no Matching facade for '%s'", alias)
	}

	return facade, nil
}

type RemotesHostsSyncFunction struct {
	db *gorm.DB
}

func (m *RemotesHostsSyncFunction) Remotes(remoteConfigs map[string]config.RemoteConnector) (*RemotesHostsFactory, error) {
	remotes := make(map[string]*remotes.RemoteFacade)

	for key, remote := range remoteConfigs {
		facade, err := m.syncRemotesWithBrain(remote)
		if err != nil {
			return nil, err
		}

		remotes[key] = facade
	}

	return &RemotesHostsFactory{
		Remotes: remotes,
	}, nil
}

func (m *RemotesHostsSyncFunction) syncRemotesWithBrain(connector config.RemoteConnector) (*remotes.RemoteFacade, error) {
	connectorURI := connector.Endpoint.URI
	if connectorURI == "" {
		switch connector.Connector {
		case config.RemoteTypeGithub:
			connectorURI = "https://github.com"
		case config.RemoteTypeGitlab:
			connectorURI = "https://gitlab.com"
		default:
			return nil, fmt.Errorf("not supported type %s", connector.Connector)
		}
	}

	model := store.RemoteRepoHostModel{
		RemoteURI:  connectorURI,
		RemoteType: connector.Connector,
	}

	lookup := store.RemoteRepoHostModel{
		RemoteURI:  connectorURI,
		RemoteType: connector.Connector,
	}

	result := m.db.FirstOrCreate(&model, &lookup)
	if result.Error != nil {
		return nil, result.Error
	}

	// get supported services from Endpoint
	services, err := getServicesFromRemoteRepoHost(connector)
	if err != nil {
		return nil, err
	}

	remoteFacade := remotes.RemoteFacade{
		Lookups: services,
		Model:   &model,
		DB:      m.db,
	}

	return &remoteFacade, nil
}

func getServicesFromRemoteRepoHost(remoteCfg config.RemoteConnector) (map[config.DirectoryType]remotes.RemoteDataLookup, error) {
	var err error

	var sevices map[config.DirectoryType]remotes.RemoteDataLookup

	switch remoteCfg.Connector {
	case config.RemoteTypeGitlab:
		sevices, err = gitlab.GetServices(remoteCfg.Endpoint)
	case config.RemoteTypeGithub:
		sevices, err = github.GetServices(remoteCfg.Endpoint)
	default:
		return nil, fmt.Errorf("not supported connector type %s", remoteCfg.Connector)
	}

	if err != nil {
		return nil, err
	}

	return sevices, nil
}
