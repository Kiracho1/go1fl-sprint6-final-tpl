package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Server оборачивает http.Server и логгер
type Server struct {
	logger *log.Logger
	server *http.Server
}

// NewServer создает HTTP-сервер с зарегистрированными хендлерами
func NewServer(logger *log.Logger) *Server {
	router := http.NewServeMux()
	router.HandleFunc("/", handlers.Index)
	router.HandleFunc("/upload", handlers.Upload)

	server := &http.Server{
		Addr:	":8080",
		Handler: router,
		ErrorLog: logger,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,

	}
	return &Server{
		logger: logger,
		server: server,
	}
}

func (s *Server) Run() error {
	s.logger.Printf("Сервер запущен на http://localhost%s", s.server.Addr)
	return s.server.ListenAndServe()
}