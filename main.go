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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		srv.ServeHTTP(w, r)
	})

	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
