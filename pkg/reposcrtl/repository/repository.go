package repository

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
)

type LocalRepository string

const (
	LocalRepositoryNotExists     LocalRepository = "NotExists"
	LocalRepositoryNeedUpdate    LocalRepository = "NeedUpdate"
	LocalRepositoryInSync        LocalRepository = "InSync"
	LocalRepositoryNotMatchLocal LocalRepository = "LocalRepoNotMatch"
	LocalRepositoryDirMissMatch  LocalRepository = "LocalDirMissMatch"
)

func (r LocalRepository) String() string {
	return string(r)
}

type Element struct {
	Remote         remotes.RemoteRepoElement
	Naming         RenameCheckoutStrategy
	CloneURILookup CloneURILookup
}

func (r *Element) Status() LocalRepository {
	repositoryPath := r.Path()
	if _, err := os.Stat(repositoryPath); os.IsNotExist(err) {
		return LocalRepositoryNotExists
	}

	gitRepo, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return LocalRepositoryDirMissMatch
	}

	remoteURLs := []string{
		r.Remote.Remote.GetCloneURI(config.HTTP),
		r.Remote.Remote.GetCloneURI(config.SSH),
	}

	matchRepo, err := isRepoMatchWithExpected(gitRepo, remoteURLs...)

	if err != nil {
		return LocalRepositoryDirMissMatch
	}

	if !matchRepo {
		return LocalRepositoryNotMatchLocal
	}

	return LocalRepositoryInSync
}

func (r *Element) Path() string {
	return r.Naming.CheckoutPath(r.Remote.Remote)
}

func (r *Element) GetCloneURI() string {
	return r.CloneURILookup.GetCloneURI(r.Remote.Remote)
}

type CloneURILookup struct {
	Protocol config.GitAccessProtocol
}

func (l *CloneURILookup) GetCloneURI(repo remotes.RemoteRepoInfo) string {
	return repo.GetCloneURI(l.Protocol)
}

func isRepoMatchWithExpected(repo *git.Repository, cloneURI ...string) (bool, error) {
	remotes, err := repo.Config()
	if err != nil {
		return false, err
	}

	originRemotes := remotes.Remotes["origin"]

	if originRemotes != nil {
		originURLs := originRemotes.URLs

		for _, expectedURL := range cloneURI {
			if utils.Contains(originURLs, expectedURL) {
				return true, nil
			}
		}
	}

	return false, nil
}
