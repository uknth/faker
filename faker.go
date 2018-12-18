package faker

import (
	"log"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Faker Fakes the request
type Faker struct {
	host string
	port string

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

func (f *Faker) run(server *Server) {
	err := server.Open()
	if err != nil {
		log.Fatal(err)
	}
}

// Close shuts down faker
func (f *Faker) Close() error {
	return f.server.Close()
}

// Open starts up faker
func (f *Faker) Open() error {
	var err error

	go f.run(f.server)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println("Recieved SIGINT, shutting down server")

	err = f.Close()
	if err != nil {
		log.Fatal("Error Closing Server: ", err.Error())
	}

	os.Exit(0)
	return err
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
	err = NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Error Initializing Faker")
	}

	faker.host = viper.GetString("server.host")
	faker.port = viper.GetString("server.port")

	server = NewServer(faker.host, faker.port)
	faker.server = server

	faker.extractHandlers()
	faker.embedHandlers()

	return faker, err
}
