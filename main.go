package main

import (
	"api5back/src/database"
	"api5back/src/server"
	"fmt"
	"net/http"
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
		if r.Method == "OPTIONS" {
			return
		}
		srv.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
