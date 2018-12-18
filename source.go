package faker

import "github.com/pkg/errors"

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

// NewSource returns text source
func NewSource(kind string, args map[string]string) (Source, error) {
	switch kind {
	case "text":
		return NewTextSource(args)
	default:
		return nil, errors.Wrap(
			errors.New("Uknown kind of source"),
			"Kind: ["+kind+"]",
		)
	}
}
