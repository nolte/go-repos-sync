package repository

import (
	"fmt"
	"path"
	"strings"

	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/remotes"
)

type RenameCheckoutStrategy struct {
	Cfg      config.Checkout
	BasePath string
}

func (r RenameCheckoutStrategy) CheckoutPath(repo remotes.RemoteRepoInfo) string {
	checkoutPath := r.BasePath
	projectName := repo.Name()

	if r.Cfg.Path != "" {
		checkoutPath = path.Join(checkoutPath, r.Cfg.Path)
	}

	if r.Cfg.KeepsGroupPath {
		var projectPath string

		projectPath = repo.Path()

		if !r.Cfg.KeepsBaseGroupNameInPath {
			// remove base dir from path
			i := strings.Index(projectPath, "/")
			if i > -1 {
				chars := projectPath[:i]
				projectPath = strings.Replace(projectPath, chars, "", 1)
			}
		}

		checkoutPath = path.Join(checkoutPath, projectPath)
	}

	if r.Cfg.FsPrefix != "" {
		projectName = fmt.Sprintf("%s%s", r.Cfg.FsPrefix, projectName)
	}

	checkoutPath = path.Join(checkoutPath, projectName)

	return checkoutPath
}
