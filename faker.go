package faker

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Faker Fakes the request
type Faker struct {
	handlers []Handler
	server   *Server
}

func (f *Faker) extractHandlers() {
	handlerConfigs := viper.Get("handlers").([]interface{})
	for _, hc := range handlerConfigs {
		if hc != nil {
			h, err := NewHandler(hc)
			if err != nil {
				log.Fatal(err)
			}

			f.handlers = append(f.handlers, h)
		}
	}
}

func (f *Faker) embedHandlers() {
	for _, h := range f.handlers {
		f.server.Handle(
			h.Path(),
			h.HandlerFunc(),
			h.Methods(),
			h.MustParams(),
		)
	}
}

func (f *Faker) callback(event fsnotify.Event) {
	// Update the config & restart the server
	f.extractHandlers()
	f.embedHandlers()

	err := f.server.Restart()
	if err != nil {
		log.Fatal(err)
	}
}

func (f *Faker) Open() error {
	return f.server.Open()
}

// NewFaker returns new object of faker
func NewFaker() (*Faker, error) {
	var (
		server *Server
		err    error
	)

	faker := &Faker{
		handlers: make([]Handler, 0),
	}

	// Parse Config
	err = NewConfig(faker.callback)
	if err != nil {
		return nil, errors.Wrap(err, "Error Initializing Faker")
	}

	server = NewServer(
		viper.GetString("server.host"),
		viper.GetString("server.port"),
	)

	faker.server = server

	faker.extractHandlers()
	faker.embedHandlers()

	return faker, err
}
