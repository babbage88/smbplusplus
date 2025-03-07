package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// ListOnlyFiles lists only the files in the specified rootPath
func ListOnlyFiles(rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking dir: %v", err)
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("error getting file info: %v", err)
		}

		// Calculate the relative path
		relativePath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("error calculating relative path: %v", err)
		}

		// Append file info
		files = append(files, FileInfo{
			FullName:     path,
			RelativeName: relativePath,
			Size:         info.Size(),
			IsDir:        false,
		})
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
