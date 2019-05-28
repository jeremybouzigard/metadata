package metadata

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jeremybouzigard/metadata"
	"github.com/jeremybouzigard/metadata/pkg/parser"
)

// Service implements metadata.Service.
type Service struct{}

// Metadata parses an audio file and returns metadata.
func (s Service) Metadata(path string) (metadata.Metadata, error) {
	var m metadata.Metadata
	var err error

	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".m4a":
		m, err = parser.M4A(path)
	default:
		return m, fmt.Errorf(metadata.UnsupportedContainerError, ext)
	}

	if err != nil {
		return m, fmt.Errorf(metadata.ParseError, path)
	}
	return m, nil
}
