package main

import (
	"fmt"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fileServer)


	http.HandleFunc("/api/donate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "success", "message": "Backend connected successfully!"}`)
	})

	fmt.Println("Server running at: http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
