package main

import (
	"cmp"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/cmd/common"
)

const serviceName = "API"
const maxBytes = 512000

type application struct {
	name     string
	port     string
	env      string
	dbURL    string
	dbAuth   string
	infoLog  *log.Logger
	errorLog *log.Logger
}

func newApplication() *application {
	return &application{
		name:     serviceName,
		port:     cmp.Or(os.Getenv("PORT"), "4000"),
		env:      cmp.Or(os.Getenv("ENV"), "development"),
		dbURL:    cmp.Or(os.Getenv("DB_URL"), ""),
		dbAuth:   cmp.Or(os.Getenv("DB_AUTH"), ""),
		infoLog:  log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func newServer(port string) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           common.Routes(),
		IdleTimeout:       20 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		MaxHeaderBytes:    maxBytes,
	}
}

func shutdownService(srv *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func newStopChannel() chan os.Signal {
	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return stopChan
}
