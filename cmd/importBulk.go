package cmd

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/syncbuilder"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	bulkConfigs string
)

func init() {
	importBulkCmd.PersistentFlags().StringVarP(&bulkConfigs, "bulkConfig", "b", "", "(optional) Overwrite, bulk configs from config file.")

	err := viper.BindPFlag("settings.bulkElements", importBulkCmd.PersistentFlags().Lookup("bulkConfig"))
	if err != nil {
		log.WithError(err).Panic("Fail to bind bulkConfig")
	}
}

var importBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "import a Set Of projects",
	Long: `Helps to keep a set of remote Projects in sync,
with Local FileSystem, by ConfigFile.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// add validation for alias and elementType combination
		return persistentPreRunEFunction(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := getSyncManagerConfig()
		log.Debugf("Import A Bulk of projects '%s'", conf.Settings.BulkElements)

		db := store.NewBrainDB(conf.Settings.Store)
		remotesFactoryFunction := syncbuilder.NewRemotesHostsSyncFunction(db)

		bulkImportCfgs, err := config.NewConfigFromFiles(conf.Settings.BulkElements)
		if err != nil {
			return err
		}

		selector := syncbuilder.ProjectBulkSelectProcess{
			RemotesFactory: remotesFactoryFunction,
			Cfg:            conf,
			BulkImportCfgs: bulkImportCfgs,
		}
		return importCommand(db, selector)
	},
}
