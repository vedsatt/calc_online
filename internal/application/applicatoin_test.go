package application_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vedsatt/calc_online/internal/application"
	"github.com/vedsatt/calc_online/pkg/calculator"
)

type ResultResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestApplication(t *testing.T) {
	testCasesSuccess := []struct {
		name        string
		expression  []byte
		expectedRes ResultResponse
		status      int
	}{
		{
			name:        "simple",
			expression:  []byte(`{"expression":"4 + 2"}`),
			expectedRes: ResultResponse{Result: 6},
			status:      http.StatusOK,
		},
		{
			name:        "priority",
			expression:  []byte(`{"expression":"( 2 + 2 ) * 2"}`),
			expectedRes: ResultResponse{Result: 8},
			status:      http.StatusOK,
		},
	}

	for _, TestCase := range testCasesSuccess {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(TestCase.expression))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		application.СalcHandler(w, req)
		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var actualRes ResultResponse
		err = json.Unmarshal(data, &actualRes)
		if err != nil {
			t.Fatal(err)
		}

		if TestCase.expectedRes != actualRes {
			t.Fatalf("Test: %s, Expected result: %v, but got: %v", TestCase.name, data, TestCase.expectedRes)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Test: %s, Expected status: %d, but got: %d", TestCase.name, http.StatusOK, res.StatusCode)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)

	w := httptest.NewRecorder()
	application.СalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var actualErr ErrorResponse
	err = json.Unmarshal(data, &actualErr)
	if err != nil {
		t.Fatal(err)
	}
	expectedErr := ErrorResponse{Error: "invalid request method"}
	if expectedErr != actualErr {
		t.Fatalf("Expected error: %s, but got: %s", expectedErr, actualErr)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status: %d, but got: %d", http.StatusMethodNotAllowed, res.StatusCode)
	}

	testCasesFail := []struct {
		name        string
		expression  []byte
		expectedErr ErrorResponse
		status      int
	}{
		{
			name:        "invalid body",
			expression:  []byte(`aaa`),
			expectedErr: ErrorResponse{Error: "invalid request body"},
			status:      http.StatusMethodNotAllowed,
		},
		{
			name:        "wrong character",
			expression:  []byte(`{"expression":"4 + a"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrWrongCharacter.Error()},
			status:      http.StatusBadRequest,
		},
		{
			name:        "empty brackets",
			expression:  []byte(`{"expression":"()"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrEmptyBrackets.Error()},
			status:      http.StatusBadRequest,
		},
		{
			name:        "division by zero",
			expression:  []byte(`{"expression":"2/(4 - 4)"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrDivisionByZero.Error()},
			status:      http.StatusBadRequest,
		},
		{
			name:        "bracket is not closed",
			expression:  []byte(`{"expression":"(4 + 2"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrNotClosedBracket.Error()},
			status:      http.StatusBadRequest,
		},
		{
			name:        "merger operators",
			expression:  []byte(`{"expression":"4 +* 2"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrMergedOperators.Error()},
			status:      http.StatusBadRequest,
		},
	}

	for _, TestCase := range testCasesFail {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(TestCase.expression))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		application.СalcHandler(w, req)
		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var actualErr ErrorResponse
		err = json.Unmarshal(data, &actualErr)
		if err != nil {
			t.Fatal(err)
		}
		if TestCase.expectedErr != actualErr {
			t.Fatalf("Expected error: %s, but got: %s", TestCase.expectedErr, data)
		}
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status: %d, but got: %d", http.StatusBadRequest, res.StatusCode)
		}
	}
}
