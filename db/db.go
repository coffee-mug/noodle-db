package db

import (
	"os"
	"log"
	"encoding/gob"
	"errors"
)

type Recipe struct {
	Name []byte
	Instructions []byte
}

func Add(name, instructions []byte) error {
	if (len(name) == 0 || len(instructions) == 0){
		return errors.New("Name and instructions are required for a new recipe")
	}
	r := Recipe{Name: name, Instructions: instructions}

	// Store it to file in binary format
	f, err := os.OpenFile("lol.db", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	enc := gob.NewEncoder(f)

	enc.Encode(r)

	return nil
}

func Load() Recipe {
	// Store it to file in binary format
	f, err := os.OpenFile("lol.db", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	enc := gob.NewDecoder(f)
	r := Recipe{}

	enc.Decode(&r)

	return r

}