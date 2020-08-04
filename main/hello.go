package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v2"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("/tmp/temp"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client := &badgerClient{db}
	client.printAllData()
	client.write("okabe",123)
	client.printAllData()
}