package store

import (
	"gorm.io/gorm"
)

type RemoteRepoModel struct {
	gorm.Model
	RemoteHostRefer uint
	Remote          RemoteRepoHostModel `gorm:"foreignKey:RemoteHostRefer"`
	RemoteLookupID  string
}
