package main

import (
	"encoding/binary"
	"fmt"
	"github.com/dgraph-io/badger/v2"
)

type badgerClient struct {
	*badger.DB
}

func (client *badgerClient) printAllData() {
	// Your code here…
	client.View(func(txn *badger.Txn) error {
		// Your code here…
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				num := int64(binary.LittleEndian.Uint64(v))
				fmt.Printf("key=%s, value=%d\n", k, num )
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (client *badgerClient) write(key string, value int) {
	client.Update(func(txn *badger.Txn) error {
		buffer := make([]byte, 8)
		binary.LittleEndian.PutUint64(buffer, uint64(value))

		entry := badger.NewEntry([]byte(key), buffer)
		return txn.SetEntry(entry)
	})
}