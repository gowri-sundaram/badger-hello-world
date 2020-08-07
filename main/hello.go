package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v2"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	options := badger.DefaultOptions("/tmp/temp").WithNumVersionsToKeep(1000)
	db, err := badger.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	client := &badgerClient{db}
	//testBasicEtag(client)
	testReadAndThenUpdate(client)
}

// Nothing happens here. You can do whatever you want, all
// Badger is doing is maintaining versions as an incremental meta.
func testBasicEtag(client *badgerClient) {
	for i := 0; i<10; i++ {
		//client.printAllData()
		client.write("lol", i * 1000)
	}
	client.printAllData()
}

// Read here is just for the sake of it.
// Otherwise, all Badger is doing is maintaining versions as an incremental meta.
func testReadAndThenUpdate(client *badgerClient) {
	for i := 0; i<10; i++ {
		//client.printAllData()
		client.readButWriteRegardlessOfRead("go", i * 1000)
	}
	client.printAllData()
}