package main

import (
	"fmt"
	"net/http"

	"api5back/src/database"
	"api5back/src/server"
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

	srv := server.NewServer(dbClient, dwClient)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "*/*")
		if r.Method == "OPTIONS" {
			return
		}
		srv.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
