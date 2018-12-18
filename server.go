package faker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Server defines server block
type Server struct {
	*http.Server

	router *mux.Router
}

// Open starts the server
func (s *Server) Open() error {
	hn := s.Prepare()

	s.Handler = hn

	log.Println("Starting Server At:", s.Addr)

	return s.ListenAndServe()
}

// Prepare binds the server with Negroni Middlewares
func (s *Server) Prepare() http.Handler {
	ng := negroni.New()

	ng.Use(negroni.NewRecovery())
	ng.Use(negroni.NewLogger())

	ng.UseHandler(s.router)
	return ng
}

// Close shuts down the server
func (s *Server) Close() error {

	ctx, cancel := context.WithTimeout(
		context.Background(), 100*time.Second,
	)

	defer cancel()

	return s.Shutdown(ctx)
}

// Handle provides utility method to handle request
func (s *Server) Handle(
	path string,
	handler http.HandlerFunc,
	method []string,
	mustParams []Pair,
) {
	// Handle the reuquest using router
	r := s.router.HandleFunc(path, handler).Methods(method...)
	for _, pr := range mustParams {
		r.Queries(pr.Key(), pr.Value())
	}
}

// NewServer returns the server object
func NewServer(host, port string) *Server {
	if port == "" {
		port = "8080"
	}
	return &Server{
		&http.Server{
			Addr: host + ":" + port,
		},
		mux.NewRouter(),
	}
}
