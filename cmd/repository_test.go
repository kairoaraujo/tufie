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
	homeDir       string
	baseDir       string
	configFile    string
}

// Mock StorageService
type storageServiceMock struct {
	mock.Mock
	storage.Storage
}

// Mock storageService.getUserHomeDir()
func (m *storageServiceMock) GetUserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestITestSuite(t *testing.T) {
	suite.Run(t, new(ITRepositorySuite))

}

func (it *ITRepositorySuite) SetupTest() {
	it.mockedStorage = new(storageServiceMock)
	it.homeDir = filepath.Join(os.TempDir(), "ITtestTUFieUser")
	it.baseDir = filepath.Join(it.homeDir, ".tufie")
	it.configFile = filepath.Join(it.baseDir, "config.yml")

	// mock the GetUserHomeDir to return the temporary dir/user
	it.mockedStorage.On("GetUserHomeDir").Return(it.homeDir, nil)

	// clean before start
	err := os.RemoveAll(it.homeDir)
	if err != nil {
		it.FailNow(err.Error())
	}
}

func (it *ITRepositorySuite) Test_Repository() {

	type testCases struct {
		name          string
		cmdArgs       []string
		expected      string
		checkEqual    bool
		checkContains bool
	}

	// define cmd.Storage as using Mocked
	Storage = storage.TufiStorageService{StgService: it.mockedStorage}

	testTable := []testCases{
		{
			name:          "`tufie repository`: without any repository configure",
			cmdArgs:       []string{"repository"},
			expected:      "Config File \"config\" Not Found in",
			checkEqual:    false,
			checkContains: true,
		},
		{
			name:          "`tufie repository add <parameter>`: Add repo rstuf as default",
			cmdArgs:       []string{"repository", "add", "--default", "--artifact-url", "https://rstuf.org", "--metadata-url", "https://metadata.rstuf.org", "--root", "../tests/test-root.json", "--name", "rstuf"},
			expected:      "\nRepository 'rstuf' added.\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository add <parameter>`: Duplicate repo rstuf as default",
			cmdArgs: []string{"repository", "add", "--default", "--artifact-url", "https://rstuf.org", "--metadata-url", "https://metadata.rstuf.org", "--root", "../tests/test-root.json", "--name", "rstuf"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nRepository 'rstuf' already exists.\n" +
				"Maybe 'artifact repository update'?\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository add <parameter>`: Add a second repository kairo, as default",
			cmdArgs: []string{"repository", "add", "-a", "https://rstuf.kairo.dev", "-m", "https://metadata.kairo.dev", "-r", "../tests/test-root.json", "-n", "kairo"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nRepository 'kairo' added.\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository`: show de default (kairo) repository",
			cmdArgs: []string{"repository"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nRepository: kairo\n" +
				"Artifact Base URL: https://rstuf.kairo.dev\n" +
				"Metadata Base URL: https://metadata.kairo.dev\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository rstuf`: show de default rstuf repository config",
			cmdArgs: []string{"repository", "rstuf"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nRepository: rstuf\n" +
				"Artifact Base URL: https://rstuf.org\n" +
				"Metadata Base URL: https://metadata.rstuf.org\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository <invalid repository>`: show invalid repository",
			cmdArgs: []string{"repository", "InexistentRepo"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\nNo repository 'InexistentRepo'.\n\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository set kairo`: set kairo (*already*) as default repository",
			cmdArgs: []string{"repository", "set", "kairo"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nNo changes. Current default repository is 'kairo'.\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository set rstuf`: set rstuf as default repository",
			cmdArgs: []string{"repository", "set", "rstuf"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nUpdated default repository to 'rstuf'.\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:          "`tufie repository set <invalid repository>`: set an invalid repository as default",
			cmdArgs:       []string{"repository", "set", "invalidRepository"},
			expected:      "Repository 'invalidRepository' doesn't exist.",
			checkEqual:    false,
			checkContains: true,
		},
		{
			name:    "`tufie repository remove kairo`: remove current default repository",
			cmdArgs: []string{"repository", "remove", "rstuf"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\nNew default repository: 'kairo'" +
				"\n\nRepository 'rstuf' removed.\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository remove kairo`: remove current default repository",
			cmdArgs: []string{"repository", "remove", "kairo"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\n\nRepository 'kairo' removed.\n",
			checkEqual:    true,
			checkContains: false,
		},
		{
			name:    "`tufie repository`: config file, but no default repository",
			cmdArgs: []string{"repository"},
			expected: "Config file used for TUFie: " + it.configFile +
				"\nNo default repository available.\n",
			checkEqual:    true,
			checkContains: false,
		},
	}

	for _, test := range testTable {
		it.T().Log(test.name)
		output := bytes.NewBufferString("")
		TUFie.SetOut(output)
		TUFie.SetErr(output)
		TUFie.SetArgs(test.cmdArgs)
		err := TUFie.Execute()
		if err != nil {
			it.FailNow(err.Error())
		}

		actual := output.String()

		if test.checkEqual {
			it.Equal(test.expected, actual)
		}
		if test.checkContains {
			it.Contains(actual, test.expected)
		}
	}
}
