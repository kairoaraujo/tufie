package storage

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock StorageService
type storageServiceMock struct {
	mock.Mock
	Storage
}

// Test Suite: UT Storage
type UTStorageSuite struct {
	suite.Suite
	mockedStorage *storageServiceMock
	stgTest       TufiStorageService
}

func TestUTestSuite(t *testing.T) {
	suite.Run(t, new(UTStorageSuite))
}

// Mock storageService.getUserHomeDir()
func (m *storageServiceMock) GetUserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (ut *UTStorageSuite) SetupTest() {
	ut.mockedStorage = new(storageServiceMock)
	ut.stgTest = TufiStorageService{ut.mockedStorage}
}

func (ut *UTStorageSuite) TestgetUserHomeDir() {

	// Use real user setup for testing
	stgService := StorageService{}
	Storage := TufiStorageService{StgService: &stgService}
	expected, err := os.UserHomeDir()
	if err != nil {
		ut.FailNow(err.Error())
	}

	homeDir, err := Storage.StgService.GetUserHomeDir()
	ut.Nil(err)
	ut.Equal(expected, homeDir)
}

func (ut *UTStorageSuite) TestGetUserHomeDir() {

	homeDir := filepath.Join(os.TempDir(), "testTUFieUser")
	// mock the getUserHomeDir to return the temporary dir/user
	ut.mockedStorage.On("GetUserHomeDir").Return(homeDir, nil)

	baseDir, err := ut.stgTest.GetBaseDir()
	ut.Equal(baseDir, filepath.Join(homeDir, ".tufie"))
	ut.Nil(err)
}

func (ut *UTStorageSuite) TestGetUserHomeDir_Error_getUserHome() {

	// mock the getUserHomeDir to return an error
	ut.mockedStorage.On("GetUserHomeDir").Return("", errors.New("Fake permission denied"))

	baseDir, err := ut.stgTest.GetBaseDir()
	ut.Error(err)
	ut.Equal("", baseDir)
}

func (ut *UTStorageSuite) TestInitDirs() {

	// mock the getUserHomeDir to use $TEMP/testTUFieUser user
	ut.mockedStorage.On("GetUserHomeDir").Return(filepath.Join(os.TempDir(), "github.com/kairoaraujo/tufie"), nil)

	err := ut.stgTest.InitDirs()
	ut.Nil(err)
}

func (ut *UTStorageSuite) TestInitDirs_Error_GetBaseDir() {

	// It also makes GetBaseDir fail
	ut.mockedStorage.On("GetUserHomeDir").Return("", errors.New("Failed GetDir"))

	err := ut.stgTest.InitDirs()
	ut.Error(err)
}

func (ut *UTStorageSuite) TestInitDirs_Error_MkdirAll() {

	// Create a $HOME as $TEMP/tufieMkdirAllFailure
	homeDir := filepath.Join(os.TempDir(), "github.com/kairoaraujo/tufieMkdirAllFailure")
	tufieDir := filepath.Join(homeDir, ".tufie")
	err := os.MkdirAll(tufieDir, 0755)
	if err != nil {
		ut.FailNow(err.Error())
	}

	// Mock the getUserHome to return
	ut.mockedStorage.On("GetUserHomeDir").Return(homeDir, nil)

	// The InitDirs creates a $HOME/tufieMkdirAllFailure/.tufie/metadata
	// To make it get an error, we create the folder as a file
	badPath := filepath.Join(tufieDir, "metadata")
	_, err = os.Create(badPath)
	if err != nil {
		ut.FailNow(err.Error())
	}

	err = ut.stgTest.InitDirs()
	ut.Error(err)
}

func (ut *UTStorageSuite) TestMakeRepository() {

	homeDir := filepath.Join(os.TempDir(), "testTUFieUser")
	// mock to use the $TEMP/testTUFieUser user
	ut.mockedStorage.On("GetUserHomeDir").Return(homeDir, nil)
	expected := filepath.Join(homeDir, ".tufie", "metadata", "testRepository")

	err := ut.stgTest.MakeRepository("testRepository")
	ut.Nil(err)

	repoDirInfo, errStat := os.Stat(expected)
	ut.Nil(errStat)
	ut.True(repoDirInfo.IsDir())
}

func (ut *UTStorageSuite) TestMakeRepository_Error_GetBaseDir() {

	// It also makes GetBaseDir fail
	ut.mockedStorage.On("GetUserHomeDir").Return("", errors.New("Fail to retrive Home"))

	err := ut.stgTest.MakeRepository("testRepository")
	ut.Error(err)
}

func (ut *UTStorageSuite) TestMakeRepository_Error_missing_metadata_dir() {

	homeDir := filepath.Join(os.TempDir(), "testTUFieUser")
	// mock to use the $TEMP/testTUFieUser user
	ut.mockedStorage.On("GetUserHomeDir").Return(homeDir, nil)

	// removes the sub-director metadata $TEMP/testTUFieUser/.tufie/metadata
	metadataDir := filepath.Join(homeDir, ".tufie", "metadata")
	err := ut.stgTest.InitDirs()
	if err != nil {
		ut.FailNow(err.Error())
	}
	os.RemoveAll(metadataDir)
	expected := filepath.Join(metadataDir, "testRepositry_metadata")

	// tries to create the repository folder, without metadata dir
	err = ut.stgTest.MakeRepository("testRepositry_metadata")
	ut.Nil(err)

	repoDirInfo, errStat := os.Stat(expected)
	ut.Nil(errStat)
	ut.True(repoDirInfo.IsDir())
}

func (ut *UTStorageSuite) TestMakeRepository_Error_missing_metadata_dir_but_fails() {

	homeDir := filepath.Join(os.TempDir(), "testTUFieUser")
	// mock to use the $TEMP/testTUFieUser user
	ut.mockedStorage.On("GetUserHomeDir").Return(homeDir, nil)

	// removes the sub-director metadata $TEMP/testTUFieUser/.tufie/metadata
	metadataDir := filepath.Join(homeDir, ".tufie", "metadata")
	err := ut.stgTest.InitDirs()
	if err != nil {
		ut.FailNow(err.Error())
	}
	os.RemoveAll(metadataDir)
	expected := filepath.Join(metadataDir, "testRepositry_metadata")

	// tries to create the repository folder, without metadata dir
	err = ut.stgTest.MakeRepository("testRepositry_metadata")
	ut.Nil(err)

	repoDirInfo, errStat := os.Stat(expected)
	ut.Nil(errStat)
	ut.True(repoDirInfo.IsDir())
}
func (ut *UTStorageSuite) TearDownSuite() {
	tempTestDir1 := filepath.Join(os.TempDir(), "testTUFieUser")
	tempTestDir2 := filepath.Join(os.TempDir(), "github.com/kairoaraujo/tufieMkdirAllFailure")
	log.Printf("Cleaning %v | %v ", tempTestDir1, tempTestDir2)
	filepath.Clean(tempTestDir1)
	filepath.Clean(tempTestDir2)
}
