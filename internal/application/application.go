package application

import (
	"encoding/json"
	"fmt"
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

type Error struct {
	Result string `json:"error"`
}

func СalcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %v %v", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("Code: %v, Invalid request method", http.StatusMethodNotAllowed)
		e := Error{Result: "invalid request method"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	log.Printf("Expression: %v", req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("Code: %v, Invalid request body", http.StatusUnprocessableEntity)
		e := Error{Result: "invalid request body"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
		return
	}

	fmt.Println(req.Expression)
	result, err := calculator.Calc(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("Code: %v, Error: %v", http.StatusUnprocessableEntity, err)
		e := Error{Result: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
		return
	}

	log.Printf("Code: %v, Result: %v", http.StatusOK, result)
	resp := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (a *Application) RunServer() {
	http.HandleFunc("/api/v1/calculate", СalcHandler)

	log.Printf("Starting server on %v", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
