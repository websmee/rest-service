package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pg/pg/v9"

	"github.com/websmee/rest-service/app"
	"github.com/websmee/rest-service/infrastructure"
	"github.com/websmee/rest-service/interfaces"
)

func main() {
	db := pg.Connect(&pg.Options{
		User:     "rest-service",
		Password: "rest-service",
		Database: "rest-service",
	})
	defer db.Close()

	if err := infrastructure.Migrate(db, "migrations/"); err != nil {
		log.Println(err)
		return
	}

	handler := interfaces.NewHTTPHandler(
		app.NewObjectProcessor(
			infrastructure.NewLocalObjectRepository(db),
			infrastructure.NewRemoteObjectRepository("http://localhost:9010"),
		),
	)

	h := http.NewServeMux()
	h.HandleFunc("/callback", handler.HandleCallback)
	go func() {
		if err := http.ListenAndServe(":9090", h); err != nil {
			log.Println(err)
		}
	}()

	// wait for OS interrupt/terminate signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
