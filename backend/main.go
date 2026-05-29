package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	http.Handle("/", http.FileServer(http.Dir("../frontend")))


	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodPost {
        http.Error(w, "POST only", http.StatusMethodNotAllowed)
        return
    }

    var body struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Invalid request"}`)
        return
    }

    if body.Email == "" || body.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Email and password required"}`)
        return
    }

    user, err := utils.LoginUser(body.Email, body.Password)
    if err != nil {
        log.Println("LoginUser error:", err)
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprint(w, `{"status":"error","message":"Invalid email or password"}`)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "success",
        "user":   user,
    })
})

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
			body.Password,
			body.AccountType,
			"", // walletAddress — empty for now
		)
		if err != nil {
			log.Println("CreateUser error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","message":"Could not create user: %s"}`, err.Error())
			return
		}

		fmt.Fprint(w, `{"status":"success","message":"Account created"}`)
	})


	http.HandleFunc("/api/donate", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodPost {
        http.Error(w, "POST only", http.StatusMethodNotAllowed)
        return
    }

    var body struct {
        CaseID        string  `json:"case_id"`
        Amount        float64 `json:"amount"`
        PaymentMethod string  `json:"payment_method"`
    }

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Invalid request"}`)
        return
    }

    if body.CaseID == "" || body.Amount <= 0 {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Missing case ID or amount"}`)
        return
    }

    // get logged in user from localStorage — anonymous for now
    donorID := r.Header.Get("X-User-ID")

    err := utils.CreateDonation(body.CaseID, donorID, body.PaymentMethod, body.Amount)
    if err != nil {
        log.Println("CreateDonation error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, `{"status":"error","message":"Donation failed: %s"}`, err.Error())
        return
    }

    // return updated case so frontend can refresh the bar
    updated, err := utils.GetCaseByID(body.CaseID)
    if err != nil {
        fmt.Fprint(w, `{"status":"success"}`)
        return
    }

    updated["status"] = "success"
    json.NewEncoder(w).Encode(updated)
})

	http.HandleFunc("/api/cases", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodGet {
        http.Error(w, "GET only", http.StatusMethodNotAllowed)
        return
    }

    cases, err := utils.GetCases()
    if err != nil {
        log.Println("GetCases error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, `{"status":"error","message":"Could not fetch cases"}`)
        return
    }

    // return empty array instead of null if no cases
    if cases == nil {
        cases = []map[string]interface{}{}
    }

    json.NewEncoder(w).Encode(cases)
})
http.HandleFunc("/api/cases/", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method != http.MethodGet {
        http.Error(w, "GET only", http.StatusMethodNotAllowed)
        return
    }

    // extract id from /api/cases/some-uuid
    id := strings.TrimPrefix(r.URL.Path, "/api/cases/")
    if id == "" {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprint(w, `{"status":"error","message":"Missing case ID"}`)
        return
    }

    c, err := utils.GetCaseByID(id)
    if err != nil {
        log.Println("GetCaseByID error:", err)
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprint(w, `{"status":"error","message":"Case not found"}`)
        return
    }

    json.NewEncoder(w).Encode(c)
})
	log.Fatal(http.ListenAndServe(":8080", nil))
}