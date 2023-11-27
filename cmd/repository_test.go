package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/kairoaraujo/tufie/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Test Suite: UT Repository
type UTRepositorySuite struct {
	suite.Suite
	repoData RepositoryData
	config   Config
}

func TestUTestSuite(t *testing.T) {
	suite.Run(t, new(UTRepositorySuite))
}

func (ut *UTRepositorySuite) SetupTest() {
	ut.repoData = RepositoryData{
		ArtifactBaseURL: "http://download.testRepo",
		MetadataURL:     "http://metadata.testRepo",
		TrustedRoot:     "RootInbase64",
	}
	repositories := map[string]RepositoryData{
		"testRepo": ut.repoData,
	}
	ut.config = Config{
		DefaultRepository: "testRepo",
		Repositories:      repositories,
	}
}

func (ut *UTRepositorySuite) Test_printRepository() {

	config := RepositoryConfig{
		repository:  "testRepo",
		metadataURL: ut.repoData.MetadataURL,
		targetURL:   ut.repoData.MetadataURL,
		trustedRoot: ut.repoData.MetadataURL,
	}
	printRepository(&config)
}

func (ut *UTRepositorySuite) Test_getRepository() {
	expected := RepositoryConfig{
		repository:  ut.config.DefaultRepository,
		metadataURL: ut.repoData.MetadataURL,
		targetURL:   ut.repoData.ArtifactBaseURL,
		trustedRoot: ut.repoData.TrustedRoot,
	}

	result, err := getRepository("testRepo", ut.config)
	ut.Nil(err)
	ut.Equal(&expected, result)
}

func (ut *UTRepositorySuite) Test_getRepository_Error_invalid_repo() {
	result, err := getRepository("invalidRepo", ut.config)
	ut.Nil(result)
	ut.Error(err)
	ut.ErrorContains(err, "No repository 'invalidRepo'.\n")
}

// Test Suite: UT Repository
type ITRepositorySuite struct {
	suite.Suite
	mockedStorage *storageServiceMock
	stgTest       storage.TufiStorageService
	homeDir       string
}

// Mock StorageService
type storageServiceMock struct {
	mock.Mock
	storage.Storage
}

// Mock storageService.getUserHomeDir()
func (m *storageServiceMock) getUserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestITestSuite(t *testing.T) {
	suite.Run(t, new(ITRepositorySuite))
}

func (it *ITRepositorySuite) SetupTest() {
	it.mockedStorage = new(storageServiceMock)
	it.stgTest = storage.TufiStorageService{StgService: it.mockedStorage}
	it.homeDir = filepath.Join(os.TempDir(), "ITtestTUFieUser")
	err := os.RemoveAll(it.homeDir)
	if err != nil {
		it.FailNow(err.Error())
	}
}

func (it *ITRepositorySuite) Test_ExecuteAnyCommand() {

	// mock the getUserHomeDir to return the temporary dir/user
	it.mockedStorage.On("getUserHomeDir").Return(it.homeDir, nil)
	stgService := storageServiceMock{}
	Storage = storage.TufiStorageService{StgService: stgService}

	actual := new(bytes.Buffer)
	TUFie.SetOut(actual)
	TUFie.SetErr(actual)
	TUFie.SetArgs([]string{"repository", "add", "help", "--artifact-url", "https://rubygems.org", "--metadata-url", "https://metadata.rubygems.org", "--root", "../tests/test-root.json", "--name", "rubygems"})
	TUFie.Execute()

	expected := "tufie repository add"

	it.Equal(actual.String(), expected, "Repository 'rubygems' added")
}
