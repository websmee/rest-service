package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pg/pg/v9"

	"github.com/websmee/rest-service/app"
	"github.com/websmee/rest-service/infrastructure"
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

	remover := app.NewExpiredObjectsRemover(infrastructure.NewLocalObjectRepository(db))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

outer:
	for {
		select {
		case <-ticker.C:
			removed, err := remover.RemoveExpired()
			if err != nil {
				log.Println(err)
			}
			log.Printf("removed %d", removed)
		case <-c: // wait for OS interrupt/terminate signal
			break outer
		}
	}
}
