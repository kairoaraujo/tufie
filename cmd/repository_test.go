package cmd

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Test Suite: UT Storage
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
