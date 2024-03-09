package tuf

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/theupdateframework/go-tuf/v2/metadata"
	"github.com/theupdateframework/go-tuf/v2/metadata/config"
	"github.com/theupdateframework/go-tuf/v2/metadata/updater"
)

// DownloadTarget downloads the target file using Updater. The Updater refreshes the top-level metadata,
// get the target information, verifies if the target is already cached, and in case it
// is not cached, downloads the target file.
func DownloadTarget(
	localMetadataDir, target, metadataURL, targetsURL, prefixDownloadDir string,
	prefixTargetsWithHash bool,
) error {
	log := metadata.GetLogger()

	rootBytes, err := os.ReadFile(filepath.Join(localMetadataDir, "root.json"))
	if err != nil {
		return err
	}

	cfg, err := config.New(metadataURL, rootBytes) // default config
	if err != nil {
		return err
	}
	cfg.LocalMetadataDir = localMetadataDir
	cfg.LocalTargetsDir = prefixDownloadDir
	cfg.RemoteTargetsURL = targetsURL
	cfg.PrefixTargetsWithHash = prefixTargetsWithHash

	// create a new Updater instance
	up, err := updater.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Updater instance: %w", err)
	}

	// try to build the top-level metadata
	err = up.Refresh()
	if err != nil {
		return fmt.Errorf("failed to refresh trusted metadata: %w", err)
	}

	// search if the desired target is available
	targetInfo, err := up.GetTargetInfo(target)
	if err != nil {
		return fmt.Errorf("target %s not found", target)
	}

	// target is available, so let's see if the target is already present locally
	path, _, err := up.FindCachedTarget(targetInfo, "")
	if err != nil {
		return fmt.Errorf("failed while finding a cached target: %w", err)
	}
	if path != "" {
		log.Info("Target is already present", "target", target, "path", path)
	}

	// target is not present locally, so let's try to download it
	path, _, err = up.DownloadTarget(targetInfo, "", "")
	if err != nil {
		return fmt.Errorf("failed to download target file %s - %w", target, err)
	}

	log.Info("Successfully downloaded target", "target", target, "path", path)

	return nil
}

func LoadTrustedRoot(filepath string) (*metadata.Metadata[metadata.RootType], error) {
	RootMetadata, err := metadata.Root().FromFile(filepath)
	if err != nil {
		return nil, err
	}

	return RootMetadata, nil

}

// Get Root from uri, which can be http/s or file
func GetRoot(uri string) ([]byte, error) {
	var rootBytes []byte
	u, _ := url.Parse(uri)
	if u.Scheme == "http" || u.Scheme == "https" {
		response, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		rb, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		rootBytes = rb
	} else {
		RootMetadata, err := LoadTrustedRoot(uri)
		if err != nil {
			return nil, err
		}
		rb, err := RootMetadata.ToBytes(false)
		if err != nil {
			return nil, err
		}
		rootBytes = rb
	}

	return rootBytes, nil
}
