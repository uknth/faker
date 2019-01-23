package faker

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	delay     = "delay"
	badstatus = "bad_status"
)

type response struct {
	Source     string            `mapstructure:"source"`
	StatusCode int               `mapstructure:"status_code"`
	Delay      int               `mapstructure:"delay"`
	Headers    map[string]string `mapstructure:"headers"`
	Arguments  map[string]string `mapstructure:"args"`
}

func (rc *response) HandlerFunc() http.HandlerFunc {
	// create the source
	source, err := NewSource(rc.Source, rc.Arguments)
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dl := r.FormValue(delay)

		if dl != "" || rc.Delay > 0 {
			if dl != "" {
				in, err := strconv.Atoi(dl)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				fmt.Println("Sleeping:", in)
				time.Sleep(time.Duration(in) * time.Second)
			} else {
				time.Sleep(time.Duration(rc.Delay) * time.Second)
			}
		}

		bt, err := source.Response()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for k, v := range rc.Headers {
			w.Header().Set(k, v)
		}

		sc := r.FormValue(badstatus)
		if sc != "" {
			in, err := strconv.Atoi(sc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(in)
		} else {
			w.WriteHeader(rc.StatusCode)
		}

		w.Write(bt)
	}
}

type handler struct {
	Pa  string            `mapstructure:"path"`
	Me  []string          `mapstructure:"methods"`
	MP  []string          `mapstructure:"must_params"`
	MKV map[string]string `mapstructure:"must_kv_params"`
	Re  response          `mapstructure:"response"`
}

func (hc *handler) Path() string      { return hc.Pa }
func (hc *handler) Methods() []string { return hc.Me }
func (hc *handler) MustParams() []Pair {
	pairs := make([]Pair, 0)
	for _, mp := range hc.MP {
		p := NewEmptyPair(mp)
		pairs = append(pairs, p)
	}

	for k, v := range hc.MKV {
		p := NewPair(k, v)
		pairs = append(pairs, p)
	}

	return pairs
}
func (hc *handler) HandlerFunc() http.HandlerFunc { return hc.Re.HandlerFunc() }

// Handler handles the request
type Handler interface {
	Path() string
	Methods() []string
	MustParams() []Pair
	HandlerFunc() http.HandlerFunc
}

// NewHandler returns the handler by reading the configuration
func NewHandler(config interface{}) (Handler, error) {
	var (
		h   handler
		err error
	)
	err = mapstructure.Decode(config, &h)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new handler")
	}
	return &h, err
}
