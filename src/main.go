package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"mantevian.xyz/cmykarte_backend/game"
	"mantevian.xyz/cmykarte_backend/handler"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/ws", handler.HandleWs(db))

	go game.StartLookingForMatches()

	log.Printf("Gateway listening on %s", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
