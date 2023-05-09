package infrastructure

import (
	"context"
	"fmt"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ivan-salazar14/firstGoPackage/domain/loyalty/infrastructure/handlers"
	"github.com/ivan-salazar14/firstGoPackage/infrastructure/database"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"

	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

func (s *Server) newServerHttp(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

func newServer(port string, conn *database.Service) *Server {

	router := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         300, // Duración del caché de preflight (en segundos)
	})

	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(corsMiddleware.Handler)

	// Crear handlers para cada dominio
	saleHandlers := handlers.NewLoyaltyHandler(conn, router)
	//default path to be used in the health checker
	saleHandlers.Routes()
	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return &Server{s}
}

func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Info().Msgf("CMD is shutting down %s", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("could not gracefully shutdown the cmd %s", err.Error())
	}

	log.Info().Msg("CMD Stopped")
}

func (srv *Server) Start() {
	log.Info().Msg("Starting API cmd")
	fmt.Println("entro start")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Could not listen on %s rv due to %s rv", srv.Addr, err.Error())
		}
	}()

	log.Info().Msgf("CMD is ready to handle requests %s", srv.Addr)
	srv.gracefulShutdown()
}

// Start aa
func Start(port string) {
	// connection to the database.
	db, err := database.NewDynamoDB()
	if err != nil {
		fmt.Printf("entro error %s\n", err)
		return
	}
	fmt.Println("entro start")
	defer func() {
		db.DB.Config.HTTPClient.CloseIdleConnections()
	}()

	server := newServer(port, db)
	// start the server.
	server.Start()
}
