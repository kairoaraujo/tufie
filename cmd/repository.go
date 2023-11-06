package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kairoaraujo/tufie/internal/tuf"
	"github.com/kairoaraujo/tufie/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	repositoryCmd = &cobra.Command{
		Use:        "repository [REPOSITORY NAME]",
		Short:      "Manage TUF repository configurations",
		Long:       ``,
		Args:       cobra.MaximumNArgs(1),
		ArgAliases: []string{"repository"},
		Run:        showRepository,
	}

	repositoryListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all repositories",
		Long:  ``,
		Run:   listRepository,
	}

	repositorySetCmd = &cobra.Command{
		Use:        "set",
		Short:      "Set the default repository",
		Long:       ``,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"repository"},
		Run:        setRepository,
	}

	repositoryAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new repository",
		Long:  ``,
		Run:   addRepository,
	}

	repositoryRemoveCmd = &cobra.Command{
		Use:        "remove",
		Short:      "Remove a repository",
		Long:       ``,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"repository"},
		Run:        removeRepository,
	}
)

func init() {
	repositoryCmd.AddCommand(repositorySetCmd)
	repositoryCmd.AddCommand(repositoryListCmd)
	repositoryCmd.AddCommand(repositoryAddCmd)
	repositoryAddCmd.PersistentFlags().StringP("name", "n", "", "repository name")
	repositoryAddCmd.PersistentFlags().StringP("root", "r", "", "trusted Root metadata")
	repositoryAddCmd.PersistentFlags().StringP("metadata-url", "m", "", "metadata URL")
	repositoryAddCmd.PersistentFlags().StringP("artifact-url", "a", "", "content artifact base URL")
	repositoryAddCmd.Flags().BoolP("default", "d", false, "set repository as default")
	repositoryAddCmd.Flags().Bool("artifact-hash", false, "add hash prefix to artifact [default: false]")
	err := repositoryAddCmd.MarkPersistentFlagRequired("name")
	cobra.CheckErr(err)
	err = repositoryAddCmd.MarkPersistentFlagRequired("metadata-url")
	cobra.CheckErr(err)
	err = repositoryAddCmd.MarkPersistentFlagRequired("root")
	cobra.CheckErr(err)
	err = repositoryAddCmd.MarkPersistentFlagRequired("artifact-url")
	cobra.CheckErr(err)
	repositoryCmd.AddCommand(repositoryRemoveCmd)
}

var config Config

type RepositoryConfig struct {
	repository  string
	metadataURL string
	targetURL   string
	trustedRoot string
}

// Prints Reposirory Configuration
func printRepository(repository *RepositoryConfig) {
	fmt.Printf("\nRepository: %v\n", repository.repository)
	fmt.Printf("Artifact Base URL: %v\n", repository.targetURL)
	fmt.Printf("Metadata Base URL: %v\n", repository.metadataURL)
}

// Gets an specific Repository configuration from Config
func getRepository(repository string, config Config) (*RepositoryConfig, error) {
	_, ok := config.Repositories[repository]
	if ok {
		return &RepositoryConfig{
			repository:  repository,
			metadataURL: config.Repositories[repository].MetadataURL,
			targetURL:   config.Repositories[repository].ArtifactBaseURL,
			trustedRoot: config.Repositories[repository].TrustedRoot,
		}, nil
	} else {
		return nil, errors.New("No repository '" + repository + "'.\n")
	}
}

func setRepository(ccmd *cobra.Command, args []string) {
	configErr := viper.ReadInConfig()
	cobra.CheckErr(configErr)

	err := viper.Unmarshal(&config)
	cobra.CheckErr(err)

	repository := args[0]
	_, ok := config.Repositories[repository]
	if ok {
		if config.DefaultRepository == repository {
			fmt.Printf("\nNo changes. Current default repository is '%v'.\n", repository)
			os.Exit(0)
		}
		config.DefaultRepository = repository
		viper.Set("default_repository", repository)
		err := viper.WriteConfigAs(viper.ConfigFileUsed())
		cobra.CheckErr(err)
		fmt.Printf("\nUpdated default repository to '%v'.\n", repository)
	} else {
		listRepository(ccmd, []string{})
		fmt.Printf("\nRepository '%v' doesn't exist.\nUse one of repositories above.\n", repository)
	}
}

func listRepository(ccmd *cobra.Command, args []string) {
	configErr := viper.ReadInConfig()
	cobra.CheckErr(configErr)

	err := viper.Unmarshal(&config)
	cobra.CheckErr(err)
	fmt.Printf("\nDefault repository: %v\n", config.DefaultRepository)

	for k := range config.Repositories {
		r, _ := getRepository(k, config)
		printRepository(r)
	}
}

func showRepository(ccmd *cobra.Command, args []string) {

	var repository string
	if len(args) == 1 {
		repository = args[0]
	} else {
		repository = ""
	}

	configErr := viper.ReadInConfig()
	cobra.CheckErr(configErr)

	err := viper.Unmarshal(&config)
	cobra.CheckErr(err)

	if repository != "" {
		cr, err := getRepository(repository, config)
		cobra.CheckErr(err)
		printRepository(cr)
	} else {
		if config.DefaultRepository != "" {
			cr, err := getRepository(config.DefaultRepository, config)
			cobra.CheckErr(err)
			printRepository(cr)

		} else {
			fmt.Println("Default repository not configured")
		}
	}
}

// Adds a new Repository to Config
func addRepository(ccmd *cobra.Command, args []string) {
	name, _ := ccmd.Flags().GetString("name")
	metadataURL, _ := ccmd.Flags().GetString("metadata-url")
	targetURL, _ := ccmd.Flags().GetString("artifact-url")
	trustedRoot, _ := ccmd.Flags().GetString("root")
	defaultRepo, _ := ccmd.Flags().GetBool("default")
	artifactHashPrefix, _ := ccmd.Flags().GetBool("artifact-hash")

	rootBytes, err := tuf.GetRoot(trustedRoot)
	cobra.CheckErr(err)

	configErr := viper.ReadInConfig()
	if configErr != nil {
		InitConfig()
		viper.Set("default_repository", name)

	} else {
		err := viper.Unmarshal(&config)
		cobra.CheckErr(err)

		_, ok := config.Repositories[name]
		if ok {
			err := errors.New(
				"\nRepository '" + name + "' already exists.\nMaybe 'artifact repository update'?\n",
			)
			cobra.CheckErr(err)
		}
	}

	if defaultRepo {
		viper.Set("default_repository", name)
	}
	viper.Set("repositories."+name+".metadata_url", metadataURL)
	viper.Set("repositories."+name+".artifact_base_url", targetURL)
	viper.Set("repositories."+name+".trusted_root", utils.EncodeTrustedRoot(rootBytes))
	viper.Set("repositories."+name+".hash_prefix", artifactHashPrefix)
	tufBaseDir, err := Storage.GetBaseDir()
	cobra.CheckErr(err)
	writeError := viper.WriteConfigAs(filepath.Join(tufBaseDir, "config.yml"))
	cobra.CheckErr(writeError)

	fmt.Printf("\nRepository '%v' added.\n", name)
}

func removeRepository(ccmd *cobra.Command, args []string) {
	repository := args[0]
	configErr := viper.ReadInConfig()
	cobra.CheckErr(configErr)

	err := viper.Unmarshal(&config)
	cobra.CheckErr(err)

	delete(viper.Get("repositories").(map[string]interface{}), repository)
	fmt.Println(len(config.Repositories))
	if config.DefaultRepository == repository {
		if len(config.Repositories) == 1 {
			viper.Set("default_repository", "")
		} else {
			for k := range config.Repositories {
				viper.Set("default_repository", k)
				fmt.Printf("New default repository: %v\n", k)
				break
			}
		}

	}
	writeError := viper.WriteConfigAs(viper.ConfigFileUsed())
	cobra.CheckErr(writeError)

}
