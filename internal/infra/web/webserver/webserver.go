package webserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]HandlerInput
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]HandlerInput),
		WebServerPort: serverPort,
	}
}

type HandlerInput struct {
	method  string
	handler http.HandlerFunc
}

func (s *WebServer) AddHandler(method string, path string, handler http.HandlerFunc) {
	s.Handlers[path] = HandlerInput{
		method:  method,
		handler: handler,
	}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handlerInput := range s.Handlers {
		s.Router.Method(handlerInput.method, path, handlerInput.handler)
	}
	err := http.ListenAndServe(s.WebServerPort, s.Router)
	if err != nil {
		slog.Error("error listening and serve", "port", s.WebServerPort, "error", err)
	}
}
