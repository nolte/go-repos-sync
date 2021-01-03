package cmd

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/report"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var brainAllCmd = &cobra.Command{
	Use:               "all",
	Short:             "All Elements",
	PersistentPreRunE: persistentPreRunEFunction,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debugf("List all Brain Elements")
		conf := getSyncManagerConfig()
		db := store.NewBrainDB(conf.Settings.Store)
		var models []store.LocalRepoModel
		result := db.Find(&models)
		if result.Error != nil {
			return result.Error
		}
		report.BrainElements(models)
		return nil
	},
}
