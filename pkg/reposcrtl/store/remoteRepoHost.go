package store

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"gorm.io/gorm"
)

type RemoteRepoHostModel struct {
	gorm.Model
	RemoteType config.RemoteTypeConnector
	RemoteURI  string `gorm:"unique"`
}
