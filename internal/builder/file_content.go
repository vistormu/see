package builder

import (
	"os"
	"strings"

	"github.com/vistormu/go-dsa/errors"
)

type FileContent struct {
	File    *File
	Content string
	NLines  int
}

func buildFileContent(args map[string]any) (*FileContent, error) {
	// file
	filePath, ok := args["file"].(string)
	if !ok || filePath == "" {
		return nil, errors.New(NotImplemented)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.New(FileNotFound).With("file", filePath).Wrap(err)
	}

	// content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New(FileNotFound).With("file", filePath).Wrap(err)
	}
	contentStr := string(content)

	// number of lines
	nLines := strings.Count(contentStr, "\n")

	return &FileContent{
		File: &File{
			Name: fileInfo.Name(),
			Path: filePath,
			Size: fileInfo.Size(),
		},
		Content: contentStr,
		NLines:  nLines,
	}, nil
}
