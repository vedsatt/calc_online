package calculator_test

import (
	"testing"

	"github.com/vedsatt/calc_online/pkg/calculator"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2-3*(3)",
			expectedResult: -1,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculator.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:       "last is operator",
			expression: "1+1*",
		},
		{
			name:       "two operators together",
			expression: "2+2**2",
		},
		{
			name:       "opened and not closed bracket",
			expression: "(2+2",
		},
		{
			name:       "division by zero",
			expression: "4 / (2 - 2)",
		},
		{
			name:       "wrong character",
			expression: "2 + 1a",
		},
		{
			name:       "no symbol between brackets",
			expression: "(2-2)(1+3)",
		},
		{
			name:       "empty",
			expression: "",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := calculator.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s is invalid but result %f was obtained", testCase.expression, val)
			}
		})
	}
}
