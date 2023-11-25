/*
Copyright © 2023 Kairo de Araujo <kairo@dearaujo.nl>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/kairoaraujo/tufie/internal/storage"

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
	cfgFile string
	Storage storage.TufiStorageService

	TUFie = &cobra.Command{
		Use:     "tufie",
		Short:   "TUF Command Line Interface",
		Long:    `The Update Framework (TUF) Command Line Interface`,
		Version: "0.1.1",
	}
)

func Execute() {
	err := TUFie.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	stgService := storage.StorageService{}
	Storage = storage.TufiStorageService{StgService: stgService}

	cobra.OnInitialize(InitConfig)
	tufBaseDir, err := Storage.GetBaseDir()
	cobra.CheckErr(err)
	TUFie.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "", "config file (default is "+tufBaseDir+"/config.yaml)",
	)
	err = viper.BindPFlag("config", TUFie.PersistentFlags().Lookup("config"))
	cobra.CheckErr(err)
	TUFie.AddCommand(downloadCmd)
	TUFie.AddCommand(repositoryCmd)

}

func InitConfig() {

	err := Storage.InitDirs()
	cobra.CheckErr(err)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		tufBaseDir, err := Storage.GetBaseDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(tufBaseDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Config file used for tuf:", viper.ConfigFileUsed())
	}

}
