package faker

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

// Source provides an interface used to get the response
type Source interface {
	Response() ([]byte, error)
}

type textSource struct {
	content string
}

func (ts *textSource) Response() ([]byte, error) { return []byte(ts.content), nil }

// NewTextSource returns TextSource that uses the content in the config file
func NewTextSource(args map[string]string) (Source, error) {
	ct, ok := args["content"]
	if !ok {
		return nil, errors.New("Insufficient Argument, missing \"content\"")
	}
	return &textSource{ct}, nil
}

type fileSource struct {
	filepath string
}

func (fs *fileSource) Response() ([]byte, error) { return ioutil.ReadFile(fs.filepath) }

// NewFileSource returns source which uses file as input
func NewFileSource(args map[string]string) (Source, error) {
	fp, ok := args["filepath"]
	if !ok {
		return nil, errors.New("Insufficient Argument, missing\"filepath\"")
	}

	abs, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}

	return &fileSource{abs}, nil
}

// NewSource returns text source
func NewSource(kind string, args map[string]string) (Source, error) {
	switch kind {
	case "text":
		return NewTextSource(args)
	case "file":
		return NewFileSource(args)
	default:
		return nil, errors.Wrap(
			errors.New("Uknown kind of source"),
			"Kind: ["+kind+"]",
		)
	}
}
