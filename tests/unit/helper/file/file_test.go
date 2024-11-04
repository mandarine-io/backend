package file_test

import (
	"github.com/mandarine-io/Backend/internal/api/helper/file"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileUtil_GetFileNameWithoutExt(t *testing.T) {
	t.Run(
		"file with extension", func(t *testing.T) {
			fileName := "example.txt"
			expected := "example"
			result := file.GetFileNameWithoutExt(fileName)
			assert.Equal(t, expected, result)
		},
	)

	t.Run(
		"file with multiple dots", func(t *testing.T) {
			fileName := "archive.tar.gz"
			expected := "archive.tar"
			result := file.GetFileNameWithoutExt(fileName)
			assert.Equal(t, expected, result)
		},
	)

	t.Run(
		"file without extension", func(t *testing.T) {
			fileName := "README"
			expected := "README"
			result := file.GetFileNameWithoutExt(fileName)
			assert.Equal(t, expected, result)
		},
	)
}

func Test_FileUtil_GetFilesFromDir(t *testing.T) {
	t.Run(
		"directory exists with files", func(t *testing.T) {
			// Create a temporary directory for testing
			dir := t.TempDir()

			// Create some test files
			file1 := filepath.Join(dir, "file1.txt")
			file2 := filepath.Join(dir, "file2.txt")
			_ = os.WriteFile(file1, []byte("test"), 0644)
			_ = os.WriteFile(file2, []byte("test"), 0644)

			files, err := file.GetFilesFromDir(dir)
			assert.NoError(t, err)
			assert.ElementsMatch(t, []string{"file1.txt", "file2.txt"}, files)
		},
	)

	t.Run(
		"directory exists with no files", func(t *testing.T) {
			// Create a temporary directory for testing
			dir := t.TempDir()

			files, err := file.GetFilesFromDir(dir)
			assert.NoError(t, err)
			assert.Empty(t, files)
		},
	)

	t.Run(
		"directory does not exist", func(t *testing.T) {
			dir := "/path/that/does/not/exist"

			files, err := file.GetFilesFromDir(dir)
			assert.Error(t, err)
			assert.Empty(t, files)
		},
	)
}
