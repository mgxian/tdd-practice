package main

import (
	"log"
	"net/http"

	"github.com/mgxian/tdd-practice/other/players"
)

const dbFilename = "game.db.json"

func main() {
	store, close, err := players.FileSystemPlayerStoreFromFile(dbFilename)
	defer close()
	if err != nil {
		log.Fatal(err)
	}

	server := players.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
