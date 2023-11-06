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
	downloadCmd = &cobra.Command{
		Use:        "download ARTIFACT",
		Short:      "Download artifact from content url using TUF metadata repository",
		Long:       ``,
		Args:       cobra.ExactArgs(1),
		ArgAliases: []string{"artifact_path"},
		Run:        download,
	}
)

func init() {
	currentDir, _ := os.Getwd()
	downloadCmd.Flags().StringP("root", "r", "", "trusted Root metadata")
	downloadCmd.Flags().StringP("metadata-url", "m", "", "metadata URL")
	downloadCmd.Flags().StringP("artifact-url", "a", "", "content artifact base URL")
	downloadCmd.Flags().StringP("directory-prefix", "P", currentDir, "save artifact to PREFIX/..")
	downloadCmd.Flags().Bool("artifact-hash", false, "add hash prefix to artifact [default: false]")
}

func download(ccmd *cobra.Command, args []string) {

	var (
		config       Config
		error_params string
	)
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		cr          string
		targetURL   string
		prefixDir   string
		metadataURL string
		trustedRoot string
		prefixHash  bool
	)

	metadataURLFlag, _ := ccmd.Flags().GetString("metadata-url")
	targetURLFlag, _ := ccmd.Flags().GetString("artifact-url")
	trustedRootFlag, _ := ccmd.Flags().GetString("root")
	prefixDir, _ = ccmd.Flags().GetString("directory-prefix") // used only on download sub-command
	prefixTargetsWithHashFlag, _ := ccmd.Flags().GetBool("artifact-hash")
	target := args[0] // map the target argument

	// if there is a default repository load it
	if config.DefaultRepository != "" {
		cr = config.DefaultRepository
		metadataURL = config.Repositories[cr].MetadataURL
		targetURL = config.Repositories[cr].ArtifactBaseURL
		trustedRoot = config.Repositories[cr].TrustedRoot
		prefixHash = config.Repositories[cr].prefixTargetsWithHash
	}

	// Flags has priority to defined configuration file
	// if the user gives metadata URL Flag overwites it
	if metadataURLFlag != "" {
		metadataURL = metadataURLFlag
	}
	// if the user gives artifact(target) URL Flag overwites it
	if targetURLFlag != "" {
		targetURL = targetURLFlag
	}
	// if the user gives trusted Root Flag, overwrites it
	if trustedRootFlag != "" {
		// load the Root in the same format a string in base64
		rootBytes, err := tuf.GetRoot(trustedRootFlag)
		cobra.CheckErr(err)
		trustedRoot = utils.EncodeTrustedRoot(rootBytes)
	}
	// if the user gives artifact(target) URL Flag overwites it
	if prefixTargetsWithHashFlag {
		prefixHash = prefixTargetsWithHashFlag
	}

	// Check if is missing configuration
	if trustedRoot == "" {
		error_params += "--root is required when no config.\n"
	}
	if metadataURL == "" {
		error_params += "--metadata-url is required when no config.\n"
	}
	if targetURL == "" {
		error_params += "--artifact-url is required when no config.\n"
	}

	if error_params != "" {
		error_params += "Use --help for more details\n"
		err := errors.New("\n" + error_params)
		cobra.CheckErr(err)
	}

	// from metadata dir defines the repoSha name and the metadata directory
	repoSha := utils.StringSha(metadataURL)
	tufBaseDir, err := Storage.GetBaseDir()
	cobra.CheckErr(err)
	metadataDir := filepath.Join(tufBaseDir, "metadata", repoSha)
	// create the repository sha folder
	err = Storage.MakeRepository(repoSha)
	cobra.CheckErr(err)

	// Save the root
	rootMetadata := utils.DecodeTrustedRoot(trustedRoot)
	cobra.CheckErr(err)
	err = rootMetadata.ToFile(filepath.Join(metadataDir, "root.json"), true)
	cobra.CheckErr(err)

	errDownload := tuf.DownloadTarget(
		metadataDir, target, metadataURL, targetURL, prefixDir, prefixHash,
	)
	cobra.CheckErr(errDownload)

	fmt.Printf("\nArtifact %v donwload completed.\n", target)
}
