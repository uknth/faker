package faker

import (
	"encoding/json"
	"fmt"
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
	return func(w http.ResponseWriter, r *http.Request) {
		// create the source
		source, err := NewSource(rc.Source, r, rc.Arguments)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dl := r.FormValue(delay)

		if dl != "" || rc.Delay > 0 {
			if dl != "" {
				in, err := strconv.Atoi(dl)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
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

type failure struct {
	Per    int `mapstructure:"percentage"`
	Status int `mapstructure:"http_status"`
}

type handler struct {
	Pa  string            `mapstructure:"path"`
	Me  []string          `mapstructure:"methods"`
	MP  []string          `mapstructure:"must_params"`
	MKV map[string]string `mapstructure:"must_kv_params"`
	Fr  failure           `mapstructure:"failure"`
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
func (hc *handler) HandlerFunc() http.HandlerFunc {
	fn := hc.Re.HandlerFunc()

	if hc.Fr.Per == 0 {
		return fn
	}

	// hc.Fr != 0
	var (
		success = 1.00
		total   = 1.00
	)

	if !(0 < hc.Fr.Per && hc.Fr.Per < 100) {
		fmt.Println("Incorrect Failure Rate:", hc.Fr.Per)
		return fn
	}

	return func(rw http.ResponseWriter, re *http.Request) {
		// fmt.Println(
		// 	"D:", float32(success/total),
		// 	"P:", float32(hc.Fr.Per)/100.00,
		// 	"B", float32(success/total) <= float32(hc.Fr.Per)/100.00,
		// )
		if float32(success/total) <= float32(hc.Fr.Per)/100.00 {
			defer func() {
				success++
				total++
			}()

			fn(rw, re)
			return
		}

		var (
			message = map[string]interface{}{
				"error": "failure due to threshold",
				"code":  hc.Fr.Status,
			}
		)

		defer func() {
			total++
		}()

		rw.WriteHeader(hc.Fr.Status)
		err := json.NewEncoder(rw).Encode(message)
		if err != nil {
			panic(err)
		}
	}

}

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
