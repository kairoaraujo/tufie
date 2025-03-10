/*
Copyright © 2023-2025 Kairo de Araujo <kairo@dearaujo.nl>
*/
package cmd

import (
	stdlog "log"
	"os"

	"github.com/go-logr/stdr"
	"github.com/kairoaraujo/tufie/internal/storage"
	"github.com/theupdateframework/go-tuf/v2/metadata"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Repository configuration data
type RepositoryData struct {
	ArtifactBaseURL       string `mapstructure:"artifact_base_url"`
	MetadataURL           string `mapstructure:"metadata_url"`
	TrustedRoot           string `mapstructure:"trusted_root"`
	prefixTargetsWithHash bool   `mapstructure:"hash_prefix"`
}

// TUFie configuration
type Config struct {
	DefaultRepository string                    `mapstructure:"default_repository"`
	Repositories      map[string]RepositoryData `mapstructure:"repositories"`
}

var (
	cfgFile   string
	verbosity bool
	Storage   storage.TufiStorageService

	TUFie = &cobra.Command{
		Use:           "tufie",
		Short:         "TUF Command Line Interface",
		Long:          `The Update Framework (TUF) Command Line Interface`,
		Version:       "0.3.1",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() {
	stgService := storage.StorageService{}
	Storage = storage.TufiStorageService{StgService: &stgService}
	err := TUFie.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)

	TUFie.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "", "config file (default is $HOME/.tufie/config.yaml)",
	)
	TUFie.PersistentFlags().BoolVarP(&verbosity, "verbose", "v", false, "verbose output")
	err := viper.BindPFlag("config", TUFie.PersistentFlags().Lookup("config"))
	cobra.CheckErr(err)

	TUFie.AddCommand(downloadCmd)
	TUFie.AddCommand(repositoryCmd)

}

func InitConfig() {
	metadata.SetLogger(stdr.New(stdlog.New(os.Stdout, "metadata - ", stdlog.LstdFlags)))
	if verbosity {
		stdr.SetVerbosity(5)
	}

	tufBaseDir, err := Storage.GetBaseDir()
	cobra.CheckErr(err)

	err = Storage.InitDirs()
	cobra.CheckErr(err)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cobra.CheckErr(err)
		viper.AddConfigPath(tufBaseDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		TUFie.Println("Config file used for TUFie:", viper.ConfigFileUsed())
	}

}

func loadConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}
