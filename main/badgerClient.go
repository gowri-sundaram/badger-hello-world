package main

import (
	"encoding/binary"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"time"
)

type badgerClient struct {
	*badger.DB
}

func (client *badgerClient) printAllData() {
	// Your code here…
	err := client.View(func(txn *badger.Txn) error {
		// Your code here…
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		opts.AllVersions = true
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			version := item.Version()
			err := item.Value(func(v []byte) error {
				num := int64(binary.LittleEndian.Uint64(v))
				fmt.Printf("key=%s, value=%d, version=%d\n", k, num, version )
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (client *badgerClient) write(key string, value int) {
	fmt.Printf("Got write call: (%s, %d)\n", key, value)
	err := client.Update(func(txn *badger.Txn) error {
		buffer := make([]byte, 8)
		binary.LittleEndian.PutUint64(buffer, uint64(value))

		entry := badger.NewEntry([]byte(key), buffer)
		return txn.SetEntry(entry)
	})
	if err != nil {
		panic(err)
	}
}

func (client *badgerClient) readButWriteRegardlessOfRead(key string, value int, sleepTimeSeconds int) error {
	fmt.Printf("Calling function with key: %s, value: %d\n", key, value)
	err := client.Update(func(txn *badger.Txn) error {
		if item, err := txn.Get([]byte(key)); err == nil {
			res, _ := item.ValueCopy(nil)
			num := int64(binary.LittleEndian.Uint64(res))
			fmt.Printf("Read value(%s): %d\n", key, num)
		}
		time.Sleep(time.Second * time.Duration(sleepTimeSeconds))
		buffer := make([]byte, 8)
		binary.LittleEndian.PutUint64(buffer, uint64(value))

		entry := badger.NewEntry([]byte(key), buffer)
		return txn.SetEntry(entry)
	})
	return err
}
