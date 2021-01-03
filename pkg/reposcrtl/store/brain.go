package store

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewBrainDB(settings config.CheckoutStoreSettings) *gorm.DB {
	databasePath := utils.ToAbsolutPath(settings.SyncDBPath)
	log.Debugf("Using Database From '%s'", databasePath)

	loggingLevel := logger.Info

	switch settings.LogLevel {
	case "warn":
		loggingLevel = logger.Warn
	case "error":
		loggingLevel = logger.Error
	}

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		Logger: logger.Default.LogMode(loggingLevel),
	})
	if err != nil {
		log.WithError(err).Panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&RemoteRepoHostModel{})
	if err != nil {
		log.WithError(err).Panic("failed to connect database")
	}

	err = db.AutoMigrate(&RemoteRepoModel{})
	if err != nil {
		log.WithError(err).Panic("failed to connect database")
	}

	err = db.AutoMigrate(&LocalRepoModel{})
	if err != nil {
		log.WithError(err).Panic("failed to connect database")
	}

	return db
}
