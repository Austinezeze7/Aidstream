package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Aidstream/backend/utils"
)

func main() {

	// =========================
	// CONNECT DATABASE
	// =========================
	err := utils.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Server running on http://localhost:8080")

	// =========================
	// STATIC FRONTEND
	// =========================
	http.Handle("/", http.FileServer(http.Dir("../frontend")))

	// =========================
	// LOGIN ROUTE
	// =========================
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		_ = json.NewDecoder(r.Body).Decode(&data)

		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
		})
	})

	// =========================
	// REGISTER ROUTE (UPDATED FOR YOUR SCHEMA)
	// =========================
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    var body struct {
        FirstName   string `json:"firstName"`
        MiddleName  string `json:"middleName"`
        LastName    string `json:"lastName"`
        Email       string `json:"email"`
        Password    string `json:"password"`
        AccountType string `json:"accountType"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Invalid request"}`)
        return
    }

    // Basic validation
    if body.Email == "" || body.Password == "" || body.FirstName == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Missing required fields"}`)
        return
    }

    err := utils.CreateUser(
        body.FirstName,
        body.MiddleName,
        body.LastName,
        body.Email,
        body.Password,  // ⚠️ hash this before storing — see note below
        body.AccountType,
        "",             // walletAddress — empty for now
    )
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"status":"error","message":"Could not create user"}`)
        return
    }

    fmt.Fprint(w, `{"status":"success","message":"Account created"}`)
})
	// =========================
	// DONATE ROUTE
	// =========================
	http.HandleFunc("/api/donate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Backend working",
		})
	})

	// =========================
	// START SERVER
	// =========================
	log.Fatal(http.ListenAndServe(":8080", nil))
}