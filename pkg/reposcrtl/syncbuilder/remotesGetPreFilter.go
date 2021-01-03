package syncbuilder

import (
	"regexp"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
	log "github.com/sirupsen/logrus"
)

type elementsGetFilter struct {
	cfg config.DirectoryContentFilter
}

func fromFilter(cfg config.DirectoryContentFilter) *elementsGetFilter {
	return &elementsGetFilter{
		cfg: cfg,
	}
}

func (f *elementsGetFilter) keep(element remotes.RemoteRepoElement) bool {
	log.Debugf("Check Element %s", element.Remote.Name())
	keeped := f.KeepsByName(element.Remote.Name())

	if keeped {
		keeped = f.KeepByArchived(element.Remote.IsArchived())
	}

	if keeped {
		keeped = f.KeepByArchivedOnly(element.Remote.IsArchived())
	}

	return keeped
}

func hasMatches(search string, patterns []string) (bool, error) {
	hasMatches := false

	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, search)
		if err != nil {
			return false, err
		}

		if matched {
			hasMatches = true
		}
	}

	return hasMatches, nil
}

func (f *elementsGetFilter) KeepByArchivedOnly(isArchivedRepo bool) bool {
	if f.cfg.ArchivedReposOnly {
		return isArchivedRepo
	}

	return !isArchivedRepo
}

func (f *elementsGetFilter) KeepByArchived(isArchivedRepo bool) bool {
	if !isArchivedRepo {
		return true
	}

	return f.cfg.IncludesArchivedRepos
}

func (f *elementsGetFilter) KeepsByName(name string) bool {
	keeps := false

	if len(f.cfg.Includes) == 0 {
		f.cfg.Includes = append(f.cfg.Includes, ".*")
	}

	matched, _ := hasMatches(name, f.cfg.Includes)
	if matched {
		keeps = true
	}

	if keeps {
		hasExcludeMatches, err := hasMatches(name, f.cfg.Excludes)
		if err != nil {
			log.WithError(err).Panicf("Fail to extract")
		}

		if hasExcludeMatches {
			return false
		}
	}

	return keeps
}
