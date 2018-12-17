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
func NewTextSource(content string) (Source, error) { return &textSource{content}, nil }

// NewSource returns text source
func NewSource(kind string, args map[string]string) (Source, error) {
	var (
		errInitSource = "Error Init Source"
	)

	switch kind {
	case "text":
		ct, ok := args["content"]
		if !ok {
			return nil, errors.Wrap(
				errors.New("Insufficient Argument, missing \"content\""),
				errInitSource,
			)
		}
		return NewTextSource(ct)
	default:
		return nil, errors.Wrap(
			errors.New("Uknown kind of source"),
			errInitSource+": "+kind,
		)
	}
}
