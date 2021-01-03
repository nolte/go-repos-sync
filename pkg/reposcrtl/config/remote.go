package config

import (
	utils "github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
)

type RemoteTypeConnector string

func (r RemoteTypeConnector) String() string {
	return string(r)
}

const (
	RemoteTypeGithub RemoteTypeConnector = "github"
	RemoteTypeGitlab RemoteTypeConnector = "gitlab"
)

type RemoteConnector struct {
	Endpoint  RemoteEndpoint      `json:"endpoint"`
	Connector RemoteTypeConnector `json:"connector"`
}

type RemoteEndpoint struct {
	URI  string              `json:"uri"`
	Auth RemoteConnectorAuth `json:"auth"`
}

type RemoteConnectorAuth struct {
	Token string `json:"token"`
}

func (a *RemoteConnectorAuth) GetTokenValue() (string, error) {
	return utils.LookupValueByRef(a.Token)
}
