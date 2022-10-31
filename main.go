package main

import (
	"github.com/Netflix/go-env"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gocloudcamp/pkg/service"
	"gocloudcamp/pkg/transport"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type environment struct {
	PgsqlURI string `env:"POSTGRES_URI"`
	HTTPPort int    `env:"HTTP_PORT"`
	GRPCPort int    `env:"GRPC_PORT"`
}

func main() {

	var e environment
	_, err := env.UnmarshalFromEnviron(&e)
	if err != nil {
		log.Fatalf("Can't get environment variables: %v", err)
	}

	svc, err := service.NewConfigService(e.PgsqlURI)
	if err != nil {
		log.Fatalf("Can't create ConfigService: %v", err)
	}

	// Running migrations
	driver, err := postgres.WithInstance(svc.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Can't get postgres driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Can't get migration object: %v", err)
	}
	if err := m.Up(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	err = transport.StartNewHTTPServer(svc, e.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	err = transport.StartNewGRPCServer(svc, e.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to start GRPC server: %v", err)
	}

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt)
	<-sigChan
	service.Shutdown(svc)
}
