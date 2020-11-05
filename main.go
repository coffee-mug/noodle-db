package main

import (
	"github.com/coffee-mug/noodle-db/db"
	"log"
)

func main() {
	r := db.Load()
	//db.Add([]byte("gâteau au chocolat"),[]byte("1/ Mélanger les oeufs (jaune + blanc) avec le sucre et le sucre vanillé"))
	log.Printf("%s", r)
}