package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		cancel()
	}()

	db := pg.Connect(&pg.Options{
		User:     "rest-service",
		Password: "rest-service",
		Database: "rest-service",
	})
	defer db.Close()

	handler := interfaces.NewHTTPHandler(
		ctx,
		app.NewObjectProcessor(
			infrastructure.NewLocalObjectRepository(db),
			infrastructure.NewRemoteObjectRepository("http://localhost:9010"),
			100,
		),
	)

	h := http.NewServeMux()
	h.HandleFunc("/callback", handler.HandleCallback)
	if err := http.ListenAndServe(":9090", h); err != nil {
		log.Println(err)
	}
}
