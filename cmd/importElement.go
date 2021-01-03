package cmd

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/syncbuilder"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	alias, elementType, id string
)

func init() {
	importElementCmd.PersistentFlags().StringVar(&alias, "alias", config.RemoteTypeGitlab.String(), "plan the checkout, no local changes if True")
	importElementCmd.PersistentFlags().StringVar(&elementType, "elementType", config.DirectoryTypeProject.String(), "Typ of imported elements")
	importElementCmd.PersistentFlags().StringVar(&id, "id", "", "Element ID")

	err := importElementCmd.MarkPersistentFlagRequired("id")
	if err != nil {
		log.WithError(err).Panic("Fail to bind id")
	}
}

var importElementCmd = &cobra.Command{
	Use:   "element",
	Short: "import a project from remote",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// add validation for alias and elementType combination
		return persistentPreRunEFunction(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := getSyncManagerConfig()
		log.Debugf("Import project %s to %s from alias '%s' Allow Dublicate Checkout '%t'", id, conf.Settings.Basedir, alias, allowDublicate)

		db := store.NewBrainDB(conf.Settings.Store)
		remotesFactoryFunction := syncbuilder.NewRemotesHostsSyncFunction(db)

		remoteFactory, err := remotesFactoryFunction.Remotes(conf.Remotes)
		if err != nil {
			return err
		}
		facade, err := remoteFactory.GetFacade(alias)
		if err != nil {
			return err
		}

		naming := repository.RenameCheckoutStrategy{
			BasePath: destinationBase,
		}

		selector := syncbuilder.ProjectSelectProcess{
			Facade: facade,
			IDs:    []string{id},
			Naming: naming,
		}
		return importCommand(db, selector)
	},
}
