package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MHG14/aethoria_marketplace/config"
	"github.com/MHG14/aethoria_marketplace/internal/application"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/service"
	"github.com/MHG14/aethoria_marketplace/internal/infrastructure/persistence/postgres"
	"github.com/MHG14/aethoria_marketplace/internal/infrastructure/persistence/postgres/repository"
	httpserver "github.com/MHG14/aethoria_marketplace/internal/transport/http"
	"github.com/MHG14/aethoria_marketplace/internal/transport/http/handlers"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load("./config")
	if err != nil {
		log.Fatal(err)
	}

	pool, err := postgres.NewPool(ctx, createDBConnString(cfg))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repos := repository.New(pool)
	app := application.NewApp(repos, service.Services{})
	h := handlers.New(app)
	srv := httpserver.NewServer(h)

	log.Fatal(srv.Listen(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)))
}

func createDBConnString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
}
