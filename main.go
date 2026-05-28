package main

import (
	"fmt"
	"Aidstream/database"
)

func main() {
	database.Connect()
	database.RunMigrations()
	fmt.Println("AidStream backend is running...")
}
