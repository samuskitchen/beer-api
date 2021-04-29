package infrastructure

import (
	"log"

	"beer-api/infrastructure/database"
)

// Start started
func Start(port string) {

	// connection to the database.
	db := database.New()
	defer db.DB.Close()

	//Versioning the database
	err := database.VersionedDB(db)
	if err != nil {
		log.Fatal(err)
	}

	server := newServer(port, db)

	// start the server.
	server.Start()
}
