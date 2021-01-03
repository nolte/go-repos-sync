package importer

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GitSyncService interface {
	action(repository.Element) error
}

type elementImporterGitSyncService struct {
	db     *gorm.DB
	lookup *repositoryElementLookupService
}

func (i elementImporterGitSyncService) action(repo repository.Element) error {
	repositoryPath := repo.Path()
	cloneURL := repo.GetCloneURI()
	log.Debugf("Git Sync Project '%s' with Local Path %s", cloneURL, repositoryPath)

	model, err := i.lookup.action(&repo)

	if err != nil {
		return err
	}

	var gitRepo *git.Repository

	if _, err := os.Stat(repositoryPath); os.IsNotExist(err) {
		log.Debug("Target dir not exists clone fresh Project")

		gitRepo, err = git.PlainClone(repositoryPath, false, &git.CloneOptions{
			URL:      cloneURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	} else {
		log.Debug("Target dir exists clone update Project")
		// check target dir is git repo and will be match for a update
		gitRepo, err = git.PlainOpen(repositoryPath)
		if err != nil {
			return err
		}
	}

	head, err := gitRepo.Head()
	if err != nil {
		return err
	}

	model.SyncSha = head.Hash().String()

	result := i.db.Save(&model)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
