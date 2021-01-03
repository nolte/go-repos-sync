package cmd

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/importer"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/syncbuilder"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	checkutProtocol, destinationBase string
	checkout, allowDublicate         bool
)

func init() {
	importCmd.AddCommand(importElementCmd)
	importCmd.AddCommand(importBulkCmd)

	// overwrite the default from global settings
	importCmd.PersistentFlags().StringVar(&destinationBase, "localPath", "", "Target checkout base directory")

	err := viper.BindPFlag("settings.basedir", importCmd.PersistentFlags().Lookup("localPath"))
	if err != nil {
		log.WithError(err).Panic("Fail to bind localPath")
	}

	importCmd.PersistentFlags().BoolVar(
		&checkout,
		"checkout",
		false,
		"checkout to local fs, if not set, no local changes will be happens.",
	)

	importCmd.PersistentFlags().BoolVar(
		&allowDublicate,
		"allowDublicate",
		false,
		"Allow Dublicate checkout to local fs",
	)

	importCmd.PersistentFlags().StringVar(
		&checkutProtocol,
		"checkutProtocol",
		config.SSH.String(),
		"Target checkout base directory",
	)

	err = viper.BindPFlag("settings.protocol", importCmd.PersistentFlags().Lookup("checkutProtocol"))
	if err != nil {
		log.WithError(err).Panic("Fail to bind checkutProtocol")
	}
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import project",
	Long:  `Import Remote Hosted Git Repostory to local FileSystem`,
}

func importCommand(db *gorm.DB, selector syncbuilder.SelectProcess) error {
	rule := importer.BrainDublicateStaticValidationRule{
		AcceptDublicateImports: allowDublicate,
	}

	var syncService importer.GitSyncService

	if checkout {
		syncService = importer.NewElementImporterGitSyncService(db)
	}

	importing := importer.ElementsImporterProcess{
		SyncService: syncService,
		Selector:    importer.NewElementImporter(selector, rule, db),
		DB:          db,
	}

	return importing.Sync()
}
