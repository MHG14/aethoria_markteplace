package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MHG14/aethoria_marketplace/config"
	"github.com/MHG14/aethoria_marketplace/internal/application"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/service"
	"github.com/MHG14/aethoria_marketplace/internal/infrastructure/adapters/oracle"
	"github.com/MHG14/aethoria_marketplace/internal/infrastructure/adapters/scheduler"
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

	mock := oracle.NewMock()
	mock.SetPrice(1, 100)
	mock.SetPrice(2, 200)
	mock.SetPrice(3, 500)
	mock.SetPrice(4, 800)
	mock.SetPrice(5, 1200)
	mock.SetPrice(6, 9999)
	mock.SetPrice(7, 9999)

	svc := service.Services{
		Oracle: mock,
		Clock:  &realClock{},
	}

	app := application.NewApp(repos, svc)
	s := scheduler.New(app, 30*time.Second)
	s.Start(ctx)
	defer s.Stop()

	h := handlers.New(app)
	srv := httpserver.NewServer(h)

	log.Fatal(srv.Listen(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)))
}

func createDBConnString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
}

type realClock struct{}

func (c *realClock) Now() time.Time { return time.Now() }
