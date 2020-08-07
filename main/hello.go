package main

import (
	"log"
	"sync"

	badger "github.com/dgraph-io/badger/v2"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	options := badger.DefaultOptions("/tmp/temp4").WithNumVersionsToKeep(1000)
	db, err := badger.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	client := &badgerClient{db}
	//testBasicEtag(client)
	//testReadAndThenUpdate(client)
	testReadAndThenUpdateWithSleep(client)
}

// All writes here succeed, may not be in the actual order.
// No conflicts seen because there is no read operation.
func testBasicEtag(client *badgerClient) {
	var wg sync.WaitGroup
	for i := 0; i<10; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			client.write("lol :D ", x*1000)
		}(i)
	}
	wg.Wait()
	client.printAllData()
}

// Read here is just dummy read to test semantics of Read/Write transactions.
// Writes fail with conflict non deterministically.
func testReadAndThenUpdate(client *badgerClient) {
	var wg sync.WaitGroup
	for i := 0; i<10; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			err := client.readButWriteRegardlessOfRead("summer", x*1000, 0)
			if err != nil {
				log.Println(err, ":", x)
			}
		}(i)
	}
	wg.Wait()
	client.printAllData()
}

// Read here is just dummy read to test semantics of Read/Write transactions.
// So here the last write succeeds, and all other writes fail deterministically.
func testReadAndThenUpdateWithSleep(client *badgerClient) {
	var wg sync.WaitGroup
	for i := 0; i<10; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			err := client.readButWriteRegardlessOfRead("autumn", x*1000, 10 - x)
			if err != nil {
				log.Println(err, ":", x)
			}
		}(i)
	}
	wg.Wait()
	client.printAllData()
}