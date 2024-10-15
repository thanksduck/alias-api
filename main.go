package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	app "github.com/thanksduck/alias-api/App"
	db "github.com/thanksduck/alias-api/Database"
	envcheck "github.com/thanksduck/alias-api/EnvCheck"
)

func main() {
	fmt.Println("Starting the application with HTTP/2 support...")
	envcheck.Init()
	db.Connect()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "6777"
	}

	// Create a new http.Server with the app.Init() handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.Init(),
	}

	fmt.Printf("Server is running on port %s with HTTP/2 support\n", port)
	fmt.Println("It also works with nodemon")

	// Start the server
	log.Fatal(server.ListenAndServe())
}
