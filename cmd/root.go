package cmd

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
	v       string
	rootCmd = &cobra.Command{
		Use:   "go-repos-sync",
		Short: "Checkouts from config",
		Long: `Keep Remote Repositories, by existing roles, in Sync with your Local Filesystem,
can keep project in sync for different Workstation.`,
		PersistentPreRunE: persistentPreRunEFunction,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.repos-sync/config.yaml)")

	rootCmd.PersistentFlags().Bool("viper", false, "use Viper for configuration")

	rootCmd.PersistentFlags().StringVarP(
		&v,
		"verbosity",
		"v",
		log.InfoLevel.String(),
		"Log level (debug, info, warn, error, fatal, panic",
	)

	err := viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	if err != nil {
		log.WithError(err).Panic("Fail to bind id")
	}

	if err := setUpLogs(os.Stdout, v); err != nil {
		er(err)
	}

	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(brainCmd)
	rootCmd.AddCommand(versionCmd)

}

func initConfig() {
	if cfgFile != "" {
		log.Debugf("Use Config file from %s", cfgFile)
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}
		log.Debugf("Try to use Config file from Home %s, if exists", home)
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(path.Join(home, ".repos-sync"))
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}

// setUpLogs set the log output ans the log level.
func setUpLogs(out io.Writer, level string) error {
	log.SetOutput(out)

	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}

	log.SetLevel(lvl)

	return nil
}

var persistentPreRunEFunction = func(cmd *cobra.Command, args []string) error {
	if err := setUpLogs(os.Stdout, v); err != nil {
		return err
	}
	return nil
}

func getSyncManagerConfig() *config.SyncManagerConfig {
	conf := &config.SyncManagerConfig{}

	err := viper.Unmarshal(conf)
	if err != nil {
		log.WithError(err).Panicf("unable to decode into config struct, %v", err)
	}

	return conf
}
