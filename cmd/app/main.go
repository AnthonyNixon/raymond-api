package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AnthonyNixon/raymond-api/cmd/services/audd"

	"github.com/AnthonyNixon/raymond-api/cmd/handlers/identify"

	"github.com/AnthonyNixon/raymond-api/cmd/handlers/auth"

	"github.com/AnthonyNixon/raymond-api/cmd/handlers/users"

	"github.com/AnthonyNixon/raymond-api/cmd/services/token"

	"github.com/AnthonyNixon/raymond-api/cmd/services/database"

	"github.com/AnthonyNixon/raymond-api/cmd/handlers/up"

	"github.com/AnthonyNixon/raymond-api/cmd/services/router"
)

var PORT = ""

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	database.Initialize()
	token.Initialize()
	audd.Initialize()
}

func main() {
	router := router.New()

	// Add Routes
	up.AddUpV1(router)
	users.AddUsersV1(router)
	auth.AddAuthV1(router)
	identify.AddIdentifyV1(router)

	log.Printf("Running Raymond API on :%s...", PORT)
	err := router.Run(fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatal(err.Error())
	}
}
