package remotes

import (
	"fmt"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"
	"gorm.io/gorm"
)

type RemoteFacade struct {
	DB      *gorm.DB
	Model   *store.RemoteRepoHostModel
	Lookups map[config.DirectoryType]RemoteDataLookup
}

type RemoteRepoElement struct {
	Remote RemoteRepoInfo
	Model  *store.RemoteRepoModel
}

func (r *RemoteFacade) GetElements(t config.DirectoryType, ids []string) ([]RemoteRepoElement, error) {
	lookup, err := r.Get(t)
	if err != nil {
		return nil, err
	}

	var allElements []RemoteRepoElement

	for _, id := range ids {
		elements, err := lookup.Get(id)
		if err != nil {
			return nil, err
		}

		for _, ele := range elements {
			var model store.RemoteRepoModel

			lookup := store.RemoteRepoModel{
				RemoteLookupID: ele.ID(),
				Remote:         *r.Model,
			}
			// brain db for element at expected path
			result := r.DB.FirstOrCreate(&model, &lookup)
			if result.Error != nil {
				return nil, result.Error
			}

			remoteElement := RemoteRepoElement{
				Remote: ele,
				Model:  &model,
			}
			allElements = append(allElements, remoteElement)
		}
	}

	return allElements, nil
}

func (r RemoteFacade) Get(t config.DirectoryType) (RemoteDataLookup, error) {
	if r.Lookups[t] == nil {
		return nil, fmt.Errorf("no LookupService for '%s'", t)
	}

	return r.Lookups[t], nil
}

type RemoteRepoInfo interface {
	ID() string
	Name() string
	IsArchived() bool
	BrowserURL() string
	GetCloneURI(config.GitAccessProtocol) string
	Path() string
}

type RemoteDataLookup interface {
	Get(string) ([]RemoteRepoInfo, error)
}
