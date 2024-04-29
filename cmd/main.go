package main

import (
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	// load env vars
	godotenv.Load()
}

func main() {
	// create channel to listen for SIGINT/SIGTERM
	stopChan := newStopChannel()

	// create app & server
	app := newApplication()
	srv := newServer(app.port)

	// start server in own goroutine
	app.infoLog.Printf("starting service. Name: %s | Mode: %s | Port: %s\n", app.name, app.env, app.port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			app.errorLog.Println(err)
			stopChan <- syscall.SIGINT
		}
	}()

	<-stopChan

	// handle SIGINT/SIGTERM. Graceful shutdown
	app.infoLog.Printf("[%s] - signal interrupted. Shutting down...", app.name)
	if err := shutdownService(srv); err != nil {
		app.errorLog.Fatal("FATAL ERROR. FORCE STOP -", err)
	}
	app.infoLog.Printf("[%s] - service stopped. Shutdown successful", app.name)
}
