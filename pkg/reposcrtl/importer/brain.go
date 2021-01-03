package importer

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"
)

type RepositorySyncBrainDublicateElementError struct {
	original   *repository.Element
	dublicates []*store.LocalRepoModel
}

func (e *RepositorySyncBrainDublicateElementError) Error() string {
	var dublicatePaths []string

	for _, dublicate := range e.dublicates {
		dublicatePaths = append(dublicatePaths, dublicate.Path)
	}

	return fmt.Sprintf("The Project '%s' allways exists at '%s'", e.original.Remote.Remote.Name(), dublicatePaths)
}

type BrainDublicateValidator struct {
	db *gorm.DB
}

func (v *BrainDublicateValidator) validate(repo repository.Element) (*store.LocalRepoModel, *RepositorySyncBrainDublicateElementError, error) {
	lookup := store.LocalRepoModel{
		RemoteRefer: repo.Remote.Model.ID,
	}

	var models []*store.LocalRepoModel
	result := v.db.Where(&lookup).Find(&models)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Debugf("No Brain Informations exists for this remote repo '%s'", repo.Remote.Remote.Name())
		return nil, nil, nil
	} else if result.Error != nil {
		return nil, nil, result.Error
	}

	var dublicateElements []*store.LocalRepoModel

	var match *store.LocalRepoModel

	for _, model := range models {
		if model.Path != repo.Path() {
			dublicateElements = append(dublicateElements, model)
		} else {
			match = model
		}
	}

	if len(dublicateElements) > 0 {
		return match, &RepositorySyncBrainDublicateElementError{
			original:   &repo,
			dublicates: dublicateElements,
		}, nil
	}

	return match, nil, nil
}

type BrainDublicateValidationRule interface {
	handle(*RepositorySyncBrainDublicateElementError) *RepositorySyncBrainDublicateElementError
}

type BrainDublicateStaticValidationRule struct {
	AcceptDublicateImports bool
}

func (v BrainDublicateStaticValidationRule) handle(err *RepositorySyncBrainDublicateElementError) *RepositorySyncBrainDublicateElementError {
	if v.AcceptDublicateImports {
		return nil
	}

	return err
}
