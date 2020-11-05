package db

import (
	"os"
	"log"
	"encoding/gob"
	"errors"
	"github.com/dghubble/trie"
)

type DB struct {
	t *trie.PathTrie
}

func NewDB() DB {
	d := DB{}

	f, err := os.OpenFile("lol.db", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	enc := gob.NewDecoder(f)

	// TODO 
	enc.Decode(&d.t)
	d.t.Segmenter = trie.PathSegmenter

	//d.t = trie.NewPathTrie()

	return d
}

func (d *DB) Keys() (keys []string) {
	d.t.Walk(func(key string, value interface{}) error {
		// debug
		if key != "" {
			keys = append(keys, key)
		}
		return nil
	})
	return keys
}

func (d *DB) Add(name, instructions string) (updated_value bool) {
	log.Printf("%v", d.t)
	return d.t.Put(name, instructions)
}

func (d *DB) Log() {
	log.Printf("Internal memory: %v", d.t)
	log.Println("All the keys:")
	for _, key := range d.Keys() {
		log.Printf("Key: %s", key)
	}
}

func (d *DB) Commit() error {
	if (d.t == nil) {
		return errors.New("Can't commit empty trie, check that your db was initialized properly")
	}
	// Store it to file in binary format
	f, err := os.OpenFile("lol.db", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", d.t)
	enc := gob.NewEncoder(f)

	enc.Encode(d.t)
	return nil
}