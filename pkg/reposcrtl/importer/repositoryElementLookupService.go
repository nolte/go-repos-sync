package importer

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"
	"gorm.io/gorm"
)

type repositoryElementLookupService struct {
	db *gorm.DB
}

func (a *repositoryElementLookupService) action(repo *repository.Element) (*store.LocalRepoModel, error) {
	lookup := store.LocalRepoModel{
		Path: repo.Path(),
	}
	model := store.LocalRepoModel{
		Path:   repo.Path(),
		Remote: *repo.Remote.Model,
		//SyncSha: repo.SHA(),
	}
	result := a.db.FirstOrCreate(&model, &lookup)

	return &model, result.Error
}
