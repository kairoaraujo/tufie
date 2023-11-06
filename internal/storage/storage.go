package storage

import (
	"os"
	"path/filepath"
)

// Manages the tufie storage
type Storage interface {
	InitDirs() error
	getUserHomeDir() (string, error)
	GetBaseDir() (string, error)
	MakeRepository(string) error
}

// Implementation of Storage Sercice
type StorageService struct {
	Storage
}

type TufiStorageService struct {
	StgService Storage
}

func (stg StorageService) getUserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// Get TUFie base directory ($HOME/.tufie)
func (ts TufiStorageService) GetBaseDir() (string, error) {
	userDir, err := ts.StgService.getUserHomeDir()
	if err != nil {
		return "", err
	}
	baseDir := filepath.Join(userDir, ".tufie")
	return baseDir, nil
}

// Initialize directories for TUFie
// - $HOME/.tufie
// - $HOME/.tufie/metadata
func (ts TufiStorageService) InitDirs() error {
	tufieDir, err := ts.GetBaseDir()
	if err != nil {
		return err
	}

	// creates $HOME/.tufie/metadata for repository data if doesnt exist
	err = os.MkdirAll(filepath.Join(tufieDir, "metadata"), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (ts TufiStorageService) MakeRepository(repoSha string) error {
	tufieDir, err := ts.GetBaseDir()
	if err != nil {
		return err
	}
	metadataDir := filepath.Join(tufieDir, "metadata")
	metadataDirInfo, errMetadataDir := os.Stat(metadataDir)
	if errMetadataDir != nil || !metadataDirInfo.IsDir() {
		errInitDir := ts.InitDirs()
		if errInitDir != nil { // excluded (it is alerady tested by InitDirs)
			return errInitDir
		}
	}

	repoDir := filepath.Join(metadataDir, repoSha)

	errMkdir := os.MkdirAll(repoDir, 0755)
	if errMkdir != nil {
		return errMkdir
	}

	return nil
}
