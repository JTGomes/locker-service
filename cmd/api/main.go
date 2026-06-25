package main

import (
	"context"
	"errors"
	"locker-service/internal/bloq"
	"locker-service/internal/config"
	"locker-service/internal/locker"
	"locker-service/internal/platform/httpserver"
	"locker-service/internal/platform/storage/postgresql"
	"locker-service/internal/rent"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := runMigrations(cfg.PostgresConfig.DSN, cfg.MigrationsPath); err != nil {
		log.Println("error running migrations ", err)
	}

	dbPool, err := postgresql.NewPool(context.Background(), &cfg.PostgresConfig)
	if err != nil {
		log.Fatal("error creating database pool", err)
	}
	defer dbPool.Close()

	transactor := postgresql.NewTransactor(dbPool)

	// repositories
	bloqRepo := postgresql.NewBloqRepository(dbPool)
	lockerRepo := postgresql.NewLockerRepository(dbPool)
	rentRepo := postgresql.NewRentRepository(dbPool)

	//services
	bloqService := bloq.NewService(bloqRepo)
	lockerService := locker.NewService(lockerRepo)
	rentService := rent.NewService(rentRepo, lockerRepo, transactor)

	// routers
	routerDep := &httpserver.RouterDependencies{
		Bloq:   bloq.NewHandler(bloqService),
		Locker: locker.NewHandler(lockerService),
		Rent:   rent.NewHandler(rentService),
	}

	router := httpserver.NewRouter(routerDep)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("API listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatal("error starting server", err)
	case sig := <-quit:
		log.Println("signal received, shutting down ", sig.String())
	}

	//graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("error shutting down server", err)
	}
	log.Println("server shut down successfully")

}

func runMigrations(databaseURL, path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	m, err := migrate.New("file://"+abs, databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
