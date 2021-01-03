package syncbuilder

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
)

type ProjectSelectProcess struct {
	IDs    []string
	Facade *remotes.RemoteFacade
	Naming repository.RenameCheckoutStrategy
}

func (i ProjectSelectProcess) GetElements() ([]repository.Element, error) {
	remoteElements, err := i.Facade.GetElements(config.DirectoryTypeProject, i.IDs)
	if err != nil {
		return nil, err
	}

	var elements []repository.Element

	for _, element := range remoteElements {
		elements = append(elements, repository.Element{
			Remote: element,
			Naming: i.Naming,
			CloneURILookup: repository.CloneURILookup{
				Protocol: config.SSH,
			},
		})
	}

	return elements, nil
}
