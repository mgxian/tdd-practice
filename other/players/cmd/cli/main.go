package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mgxian/tdd-practice/other/players"
)

const dbFilename = "game.db.json"

func main() {
	store, close, err := players.FileSystemPlayerStoreFromFile(dbFilename)
	defer close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	players.NewCLI(store, os.Stdin).PlayPoker()
}
