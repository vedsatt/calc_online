package application

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vedsatt/calc_online/pkg/calculator"
)

const (
	port = ":8080"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %v %v", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		log.Printf("Code: %v, Invalid request method", http.StatusMethodNotAllowed)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	log.Printf("Expression: %v", req)
	if err != nil {
		log.Printf("Code: %v, Invalid request body", http.StatusBadRequest)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := calculator.Calc(req.Expression)
	if err != nil {
		log.Printf("Code: %v, Error: %v", http.StatusBadRequest, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Code: %v, Result: %v", http.StatusOK, result)
	resp := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (a *Application) RunServer() {
	http.HandleFunc("/api/v1/calculate", calcHandler)

	log.Printf("Starting server on %v", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
