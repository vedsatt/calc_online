package calculator

import (
	"strconv"
	"strings"
)

func expressionErrors(expression string) error {
	len := len(expression)
	flag := false
	start := 0
	end := 0

	for i := 0; i < len; i++ {
		curr := expression[i]
		next := byte(0)
		if i < len-1 {
			next = expression[i+1]
		}

		if curr == '(' {
			start++
		}
		if curr == ')' {
			end++
		}
		if 48 <= curr && curr <= 57 && !flag {
			flag = true
		}

		switch {
		case i == 0 && (curr == ')' || curr == '*' || curr == '+' || curr == '-' || curr == '/'):
			return ErrOperatorFirst
		case i == len-1 && (curr == '(' || curr == '*' || curr == '+' || curr == '-' || curr == '/'):
			return ErrOperatorLast
		case curr == '(' && next == ')':
			return ErrEmptyBrackets
		case curr == ')' && next == '(':
			return ErrMergedBrackets
		case (curr == '*' || curr == '+' || curr == '-' || curr == '/') && (next == '*' || next == '+' || next == '-' || next == '/'):
			return ErrMergedOperators
		case curr != ' ' && (curr < '(' || curr > '9'):
			return ErrWrongCharacter
		case len <= 2:
			return ErrInvalidExpression
		}
	}

	if start > end {
		return ErrNotClosedBracket
	} else if end > start {
		return ErrNotOpenedBracket
	}
	if !flag {
		return ErrNoOperators
	}
	return nil
}

func (s *Stack) lineToStacks(expression string) {
	var tmp string
	var len int = len([]rune(expression))

	for index, char := range expression {
		switch {
		case '0' <= char && char <= '9' || char == '.' || char == ',':
			tmp += string(char)
			if index == len-1 {
				num, _ := strconv.ParseFloat(tmp, 64)
				s.numbers = append(s.numbers, num)
				tmp = ""
			}
		case char == '(' || char == ')' || char == '*' || char == '+' || char == '-' || char == '/':
			if tmp != "" {
				num, _ := strconv.ParseFloat(tmp, 64)
				s.numbers = append(s.numbers, num)
				tmp = ""
			}
			s.operators = append(s.operators, string(char))
		}
	}
}

type Stack struct {
	numbers   []float64
	operators []string
}

type StackOperators interface {
	Push(interface{})
	Pop(string) interface{}
}

func (s *Stack) push(item interface{}) {
	switch char := item.(type) {
	case float64:
		s.numbers = append(s.numbers, char)
	case string:
		s.operators = append(s.operators, char)
	}
}

func (s *Stack) pop(StackType string) interface{} {
	switch StackType {
	case "num":
		len := len(s.numbers)
		value := s.numbers[len-1]
		s.numbers = s.numbers[:len-1]
		return value
	case "op":
		len := len(s.operators)
		value := s.operators[len-1]
		s.operators = s.operators[:len-1]
		return value
	}
	return 0
}

func operations(x, y float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return x + y, nil
	case "-":
		return x - y, nil
	case "*":
		return x * y, nil
	case "/":
		if y == 0 {
			return 0, ErrDivisionByZero
		}
		return x / y, nil
	}
	return 0, ErrUnknownOperator
}

func priority(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func Calc(expression string) (float64, error) {
	expression = strings.TrimSpace(expression)
	err := expressionErrors(expression)
	if err != nil {
		return 0, err
	}

	tempNum := make([]float64, 0)
	tempOp := make([]string, 0)
	flag := false
	bracketNum := 0.0
	bracketOp := ""
	s := Stack{}
	s.lineToStacks(expression)

	for len(s.operators) > 0 {
		op := s.pop("op").(string)

		if op == ")" {
			flag2 := true
			brackNum := make([]float64, 0)
			brackOp := make([]string, 0)

			for len(s.operators) > 0 && s.operators[len(s.operators)-1] != "(" {
				op = s.pop("op").(string)

				if priority(op) == 2 || !flag2 {
					y := s.pop("num").(float64)
					x := s.pop("num").(float64)
					result, err := operations(x, y, op)

					if err != nil {
						return 0, err
					}
					s.push(result)

				} else {
					switch op {
					case "-":
						brackOp = append(brackOp, "+")
						brackNum = append(brackNum, (s.pop("num").(float64) * (-1)))
					case "+":
						brackOp = append(brackOp, op)
						brackNum = append(brackNum, s.pop("num").(float64))
					}
				}

				if s.operators[len(s.operators)-1] == "(" && flag2 {
					s.operators = append(s.operators, brackOp...)
					s.numbers = append(s.numbers, brackNum...)
					flag2 = false
				}
			}

			s.pop("op")
			if bracketOp != "" {
				s.operators = append(s.operators, bracketOp)
				s.numbers = append(s.numbers, bracketNum)
				bracketOp = ""
				bracketNum = 0.0
			}

		} else {
			if ((priority(op) == 2 || flag) && len(s.operators) == 0) || ((priority(op) == 2 || flag) && s.operators[len(s.operators)-1] != ")") {
				y := s.pop("num").(float64)
				x := s.pop("num").(float64)
				result, err := operations(x, y, op)

				if err != nil {
					return 0, err
				}
				s.push(result)

			} else {
				switch op {
				case "-":
					tempOp = append(tempOp, "+")
					tempNum = append(tempNum, (s.pop("num").(float64) * (-1)))
				case "+":
					tempOp = append(tempOp, op)
					tempNum = append(tempNum, s.pop("num").(float64))
				}
			}
		}

		if len(s.operators) != 0 {
			if s.operators[len(s.operators)-1] == ")" {
				bracketNum = s.pop("num").(float64)
				bracketOp = op

			}
		}

		if len(s.operators) == 0 && !flag {
			s.operators = tempOp
			s.numbers = append(s.numbers, tempNum...)
			flag = true
		}
	}

	return s.numbers[0], nil
}
