package importer

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
)

type elementImporterProcess struct {
	repo repository.Element
	sync GitSyncService
}

func (i *elementImporterProcess) Sync() error {
	if i.sync != nil {
		return i.sync.action(i.repo)
	}

	return nil
}
