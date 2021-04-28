package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pg/pg/v9"

	"github.com/websmee/rest-service/app"
	"github.com/websmee/rest-service/infrastructure"
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

	remover := app.NewExpiredObjectsRemover(infrastructure.NewLocalObjectRepository(db))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

outer:
	for {
		select {
		case <-ticker.C:
			err := remover.RemoveExpired(ctx)
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			break outer
		}
	}
}
