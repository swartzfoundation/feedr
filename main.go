package main

import (
	"log"
	"log/slog"
	"net/http"
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

	frontendFS := frontend.FS()
	r.Get("/*", http.FileServer(frontendFS).ServeHTTP)
	slog.Info("Build", "Time", BuildTime)
	slog.Info("Build", "Version", Version)
	slog.Info("Starting Feedr", "Port", cfg.PORT)
	err := http.ListenAndServe(":"+cfg.PORT, r)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
