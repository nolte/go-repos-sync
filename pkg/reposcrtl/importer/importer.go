package importer

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/report"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/syncbuilder"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewElementImporter(
	selector syncbuilder.SelectProcess,
	brainValidationRule BrainDublicateValidationRule,
	db *gorm.DB,
) BrainImportElementSelectorProcess {
	return BrainImportElementSelectorProcess{
		Selector:            selector,
		BrainValidationRule: brainValidationRule,
		BrainDublicateValidator: BrainDublicateValidator{
			db: db,
		},
	}
}

type BrainImportElementSelectorProcess struct {
	Selector                syncbuilder.SelectProcess
	BrainValidationRule     BrainDublicateValidationRule
	BrainDublicateValidator BrainDublicateValidator
}

func (i *BrainImportElementSelectorProcess) GetElements() ([]repository.Element, []*RepositorySyncBrainDublicateElementError, error) {
	var allElements []repository.Element

	var brainErrors []*RepositorySyncBrainDublicateElementError

	elements, err := i.Selector.GetElements()
	if err != nil {
		return nil, nil, err
	}

	log.Debugf("Validates '%d' elements with Brain Status ", len(elements))

	for _, e := range elements {
		log.Debugf("Validate '%s' to paned checkout '%s'", e.Remote.Remote.Name(), e.Path())

		match, dublicateError, err := i.BrainDublicateValidator.validate(e)
		if err != nil {
			return nil, nil, err
		}

		brainErr := i.BrainValidationRule.handle(dublicateError)
		if brainErr == nil || match != nil {
			allElements = append(allElements, e)
		} else {
			brainErrors = append(brainErrors, brainErr)
			log.Debugf("Remove because: %s", brainErr.Error())
		}
	}

	return allElements, brainErrors, nil
}

func NewElementImporterGitSyncService(db *gorm.DB) GitSyncService {
	return elementImporterGitSyncService{
		db: db,
		lookup: &repositoryElementLookupService{
			db: db,
		},
	}
}

type ElementsImporterProcess struct {
	Selector    BrainImportElementSelectorProcess
	DB          *gorm.DB
	SyncService GitSyncService
}

func (i *ElementsImporterProcess) Sync() error {
	elements, dublicates, err := i.Selector.GetElements()
	if err != nil {
		return err
	}

	log.Debugf("Sync '%d' elements with FS, '%d' are dublicates", len(elements), len(dublicates))

	for _, element := range elements {
		elementSync := &elementImporterProcess{
			sync: i.SyncService,
			repo: element,
		}
		err := elementSync.Sync()

		if err != nil {
			if !errors.Is(err, transport.ErrEmptyRemoteRepository) {
				return err
			}
		}
	}

	report.RepositoryElement(elements)

	return nil
}
