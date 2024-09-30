package main

import (
	"api5back/src/database"
	"api5back/src/server"
	"fmt"
)

func main() {
	dbClient, err := database.Setup("DB")
	if err != nil {
		panic(fmt.Errorf("failed to setup normalized database: %v", err))
	}
	defer dbClient.Close()

	dwClient, err := database.Setup("DW")
	if err != nil {
		panic(fmt.Errorf("failed to setup data warehouse: %v", err))
	}
	defer dwClient.Close()

	server.
		NewServer(dbClient, dwClient).
		Run(":8080")
}
