package syncbuilder

import (
	"path"
	"path/filepath"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"

	log "github.com/sirupsen/logrus"
)

type ProjectBulkSelectProcess struct {
	RemotesFactory *RemotesHostsSyncFunction
	Cfg            *config.SyncManagerConfig
	BulkImportCfgs []config.Config
}

func (i ProjectBulkSelectProcess) GetElements() ([]repository.Element, error) {
	var allElements []repository.Element

	allFacades := RemotesHostsFactory{
		Remotes: make(map[string]*remotes.RemoteFacade),
	}

	remoteFactory, err := i.RemotesFactory.Remotes(i.Cfg.Remotes)
	if err != nil {
		return nil, err
	}

	for k, v := range remoteFactory.Remotes {
		allFacades.Remotes[k] = v
	}

	for _, bulkImportCfg := range i.BulkImportCfgs {
		log.Debugf("Load Bulk Config with '%d' Remotes and '%d' Directory Elements", len(bulkImportCfg.Remotes), len(bulkImportCfg.Directories))

		remoteFactory, err := i.RemotesFactory.Remotes(bulkImportCfg.Remotes)
		if err != nil {
			return nil, err
		}

		for k, v := range remoteFactory.Remotes {
			allFacades.Remotes[k] = v
		}

		for _, directory := range bulkImportCfg.Directories {
			log.Debugf("Prepare Import Elements for Directory '%s' with '%d' Elements", directory.Path, len(directory.Contents))

			for _, content := range directory.Contents {
				log.Debugf("Get Elements From Alias '%s' with  IDSs: '%s' Of Type '%s'", content.Ref.Remote, content.Ref.IDS, content.Ref.Kind)

				remoteFacade, err := allFacades.GetFacade(content.Ref.Remote)
				if err != nil {
					return nil, err
				}

				remoteElements, err := remoteFacade.GetElements(content.Ref.Kind, content.Ref.IDS)
				if err != nil {
					return nil, err
				}

				// filter the results
				filter := fromFilter(content.Ref.Filter)

				for _, e := range remoteElements {
					if filter.keep(e) {
						allElements = appendByFilter(allElements, directory, e, i.Cfg.Settings, content.Checkout)
					}
				}
			}
		}
	}

	// filter list by dublicate elements
	return allElements, nil
}

func appendByFilter(
	allElements []repository.Element,
	directory config.Directory,
	e remotes.RemoteRepoElement,
	settings config.CheckoutSettings,
	checkout config.Checkout,
) []repository.Element {
	var basePath string
	if !filepath.IsAbs(directory.Path) {
		basePath = path.Join(utils.ToAbsolutPath(settings.Basedir), directory.Path)
	} else {
		basePath = directory.Path
	}

	currentElement := repository.Element{
		Remote: e,
		Naming: repository.RenameCheckoutStrategy{
			Cfg:      checkout,
			BasePath: basePath,
		},
		CloneURILookup: repository.CloneURILookup{
			Protocol: settings.DefaultProtocol,
		},
	}

	if !isDublicateImport(allElements, currentElement) {
		allElements = append(allElements, currentElement)
	}

	return allElements
}

func isDublicateImport(allElements []repository.Element, element repository.Element) bool {
	for _, existingElement := range allElements {
		if existingElement.Remote.Model.ID == element.Remote.Model.ID {
			return true
		}
	}

	return false
}
