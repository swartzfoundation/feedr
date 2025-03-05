package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swartzfoundation/feedr/frontend"
	"github.com/swartzfoundation/feedr/pkg/config"
)

var BuildTime string // seconds since 1970-01-01 00:00:00 UTC
var Version = "development"

func main() {
	var srv http.Server
	cfg := config.Load()

	if cfg.DEBUG {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	maxAge := 12 * time.Hour
	corz := cors.Options{
		AllowedOrigins:   cfg.ALLOWED_ORIGINS,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(maxAge),
	}
	r.Use(cors.Handler(corz))

	r.Get("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/*", frontend.HandlerFn())
	slog.Info("Build", "Time", BuildTime)
	slog.Info("Build", "Version", Version)
	slog.Info("Starting Feedr", "Port", cfg.PORT)

	srv = http.Server{
		Addr:    ":" + cfg.PORT,
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			slog.Error("HTTP server Shutdown", "error", err)
		}
		close(idleConnsClosed)
	}()

	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("http server listen", "error", err.Error())
		os.Exit(1)
	}

	<-idleConnsClosed
}
