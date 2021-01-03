package config

import (
	"fmt"

	"github.com/ghodss/yaml"
)

const (
	knownAPIVersion = "reposcrtl.nolte.github.com/v1alpha1"
	knownKind       = "Config"
)

type SyncManagerConfig struct {
	APIVersion string                     `json:"apiVersion"`
	Kind       string                     `json:"kind"`
	Remotes    map[string]RemoteConnector `json:"remotes,omitempty"`
	Settings   CheckoutSettings           `json:"settings,omitempty"`
}

type CheckoutSettings struct {
	Basedir         string                `json:"basedir"`
	DefaultProtocol GitAccessProtocol     `json:"protocol"`
	Store           CheckoutStoreSettings `json:"store"`
	BulkElements    []string              `json:"bulkElements,omitempty"`
}
type CheckoutStoreSettings struct {
	SyncDBPath string `json:"syncDBPath"`
	LogLevel   string `json:"logLevel"`
}

type Config struct {
	APIVersion  string                     `json:"apiVersion"`
	Kind        string                     `json:"kind"`
	Remotes     map[string]RemoteConnector `json:"remotes,omitempty"`
	Directories []Directory                `json:"directories,omitempty"`
}

func NewConfigFromFiles(paths []string) ([]Config, error) {
	var configs []Config

	err := parseResources(paths, func(docBytes []byte) error {
		var res resource
		err := yaml.Unmarshal(docBytes, &res)
		if err != nil {
			return fmt.Errorf("unmarshaling doc: %s", err)
		}
		switch {
		case res.APIVersion == knownAPIVersion && res.Kind == knownKind:
			config, err := NewConfigFromBytes(docBytes)
			if err != nil {
				return fmt.Errorf("unmarshaling config: %s", err)
			}
			configs = append(configs, config)

		default:
			return fmt.Errorf("unknown apiVersion '%s' or kind '%s' for resource",
				res.APIVersion, res.Kind)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func NewConfigFromBytes(bs []byte) (Config, error) {
	var config Config

	err := yaml.Unmarshal(bs, &config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshaling config: %s", err)
	}

	return config, nil
}
