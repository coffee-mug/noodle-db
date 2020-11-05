package main

import (
	"github.com/coffee-mug/noodle-db/db"
)

func main() {
	DB := db.NewDB()	
	DB.Log()
}