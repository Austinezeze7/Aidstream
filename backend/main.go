package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"Aidstream/backend/utils"
)

func main() {

	err := utils.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fileServer := http.FileServer(http.Dir("../frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/html/signin.html")
	})

	http.HandleFunc("/api/donate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		response := map[string]string{
			"status":  "success",
			"message": "Backend connected successfully!",
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Server running at: http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}