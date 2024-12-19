package calculator

import (
	"errors"
)

var (
	ErrOperatorFirst     = errors.New("the first character is the operator")
	ErrOperatorLast      = errors.New("the last character is the operator")
	ErrEmptyBrackets     = errors.New("empty brackets")
	ErrMergedBrackets    = errors.New("no symbol between brackets")
	ErrMergedOperators   = errors.New("the two operators are next to each other")
	ErrWrongCharacter    = errors.New("the wrong character was found")
	ErrInvalidExpression = errors.New("invalid expression")
	ErrNotOpenedBracket  = errors.New("the bracket is not closed")
	ErrNotClosedBracket  = errors.New("the bracket is not open")
	ErrNoOperators       = errors.New("operators not found")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrUnknownOperator   = errors.New("unknown operator")
)
