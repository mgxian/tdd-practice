package players

import (
	"log"
	"net/http"
	"os"
)

const dbFilename = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFilename, err)
	}

	store, err := NewFileSystemStore(db)
	if err != nil {
		log.Fatalf("problem creating file system store: %v", err)
	}
	server := newPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
