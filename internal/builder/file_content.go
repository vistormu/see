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

func buildFileContent(args Args) (*FileContent, error) {
	// content
	content, err := os.ReadFile(args.Element)
	if err != nil {
		return nil, errors.New(FileNotFound).With("file", args.Element).Wrap(err)
	}
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "\t", "    ")
	contentStr = strings.TrimRight(contentStr, "\n")

	// info
	fileInfo, err := os.Stat(args.Element)
	if err != nil {
		return nil, errors.New(FileInfo).With("path", args.Element).Wrap(err)
	}

	// number of lines
	nLines := strings.Count(contentStr, "\n") + 1

	return &FileContent{
		File: &File{
			Name: fileInfo.Name(),
			Path: args.Element,
			Size: fileInfo.Size(),
		},
		Content: contentStr,
		NLines:  nLines,
	}, nil
}
