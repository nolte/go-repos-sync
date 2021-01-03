package store

import (
	"gorm.io/gorm"
)

type LocalRepoStatus string

const (
	LocalStatusPlaned LocalRepoStatus = "planed"
	LocalStatusExists LocalRepoStatus = "exists"
)

type LocalRepoModel struct {
	gorm.Model
	RemoteRefer uint
	Remote      RemoteRepoModel `gorm:"foreignKey:RemoteRefer"`
	Path        string          `gorm:"unique"`
	SyncSha     string
	Status      LocalRepoStatus
}
