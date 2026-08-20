package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	Folio "github.com/CephandriusMaxtori/Folio"
	"github.com/CephandriusMaxtori/Folio/internal/config"
	"github.com/CephandriusMaxtori/Folio/internal/database"
	"github.com/CephandriusMaxtori/Folio/internal/handlers"
	"github.com/CephandriusMaxtori/Folio/internal/middleware"
	"github.com/CephandriusMaxtori/Folio/internal/services"
	"github.com/go-chi/chi/v5"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "Show version and exit")
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *showVersion {
		fmt.Printf("folio %s (built %s)\n", version, buildTime)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Open(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	svc := services.New(db)
	h := handlers.New(svc, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.CORS)

	h.RegisterRoutes(r)

	r.Handle("/*", http.FileServer(http.FS(Folio.WebFiles)))

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Folio %s starting on %s", version, addr)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	server.Close()
}
