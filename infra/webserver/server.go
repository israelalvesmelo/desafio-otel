package webserver

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) AddMiddleware(m func(http.Handler) http.Handler) {
	s.Router.Use(m)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	log.Println("Starting web server on port", s.WebServerPort)

	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		log.Println("server error", "error", err)
		os.Exit(1)
	}
}
