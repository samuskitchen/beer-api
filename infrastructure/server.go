package infrastructure

import (
	"beer-api/infrastructure/database"
	"beer-api/infrastructure/middleware"
	"context"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Server is a base Server configuration.
type Server struct {
	*http.Server
}

// ServeHTTP implements the http.Handler interface for the server type.
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Handler.ServeHTTP(w, r)
}

// NewServerTest initialized a Routes Server with configuration for tests.
func NewServerTest(port string, conn *database.Data) *Server {
	return newServer(port, conn)
}

// newServer initialized a Routes Server with configuration.
func newServer(port string, conn *database.Data) *Server {

	router := chi.NewRouter()

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(middleware.CORSMiddleware)

	router.Mount("/api/v1", Routes(conn))

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{s}
}

// Start the server.
func (srv *Server) Start() {
	log.Println("starting API cmd")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s rv due to %s rv", srv.Addr, err.Error())
		}
	}()
	log.Printf("cmd is ready to handle requests %s", srv.Addr)
	srv.gracefulShutdown()
}

func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Printf("cmd is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the cmd %s", err.Error())
	}
	log.Printf("cmd stopped")
}
