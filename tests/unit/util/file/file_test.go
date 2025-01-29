package file

import (
	"github.com/mandarine-io/backend/internal/util/file"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"os"
	"path/filepath"
	"testing"
)

type FileUtilSuite struct {
	suite.Suite
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(FileUtilSuite))
}

func (s *FileUtilSuite) Test_GetFileNameWithoutExt_FileWithExtension(t provider.T) {
	t.Title("GetFileNameWithoutExt - file with extension")
	t.Severity(allure.NORMAL)
	t.Tags("positive")
	t.Parallel()

	fileName := "example.txt"
	expected := "example"
	result := file.GetFileNameWithoutExt(fileName)
	t.Require().Equal(expected, result)
}

func (s *FileUtilSuite) Test_GetFileNameWithoutExt_FileWithMultipleDots(t provider.T) {
	t.Title("GetFileNameWithoutExt - file with multiple dots")
	t.Severity(allure.NORMAL)
	t.Tags("positive")
	t.Parallel()

	fileName := "archive.tar.gz"
	expected := "archive.tar"
	result := file.GetFileNameWithoutExt(fileName)
	t.Require().Equal(expected, result)
}

func (s *FileUtilSuite) Test_GetFileNameWithoutExt_FileWithoutExtension(t provider.T) {
	t.Title("GetFileNameWithoutExt - file without extension")
	t.Severity(allure.NORMAL)
	t.Tags("positive")
	t.Parallel()

	fileName := "README"
	expected := "README"
	result := file.GetFileNameWithoutExt(fileName)
	t.Require().Equal(expected, result)
}

func (s *FileUtilSuite) Test_GetFilesFromDir_DirectoryExistsWithFiles(t provider.T) {
	t.Title("GetFilesFromDir - directory exists with files")
	t.Severity(allure.NORMAL)
	t.Tags("positive")
	t.Parallel()

	// Create a temporary directory for testing
	dir := t.TempDir()

	// Create some test files
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")
	_ = os.WriteFile(file1, []byte("test"), 0644)
	_ = os.WriteFile(file2, []byte("test"), 0644)

	files, err := file.GetFilesFromDir(dir)
	t.Require().NoError(err)
	t.Require().ElementsMatch([]string{"file1.txt", "file2.txt"}, files)
}

func (s *FileUtilSuite) Test_GetFilesFromDir_DirectoryExistsWithNoFiles(t provider.T) {
	t.Title("GetFilesFromDir - directory exists with no files")
	t.Severity(allure.NORMAL)
	t.Tags("positive")
	t.Parallel()

	// Create a temporary directory for testing
	dir := t.TempDir()

	files, err := file.GetFilesFromDir(dir)
	t.Require().NoError(err)
	t.Require().Empty(files)
}

func (s *FileUtilSuite) Test_GetFilesFromDir_DirectoryNotExists(t provider.T) {
	t.Title("GetFilesFromDir - directory not exists")
	t.Severity(allure.CRITICAL)
	t.Tags("negative")
	t.Parallel()

	dir := "/path/that/does/not/exist"

	files, err := file.GetFilesFromDir(dir)
	t.Require().Error(err)
	t.Require().Empty(files)
}
