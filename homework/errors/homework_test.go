package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}

	b := &strings.Builder{}
	for _, err := range e.errs {
		b.Write([]byte("\t* "))
		b.WriteString(err.Error())
	}

	return fmt.Sprintf("%d errors occured:\n%s\n", len(e.errs), b.String())
}

func (e *MultiError) add(errs ...error) *MultiError {
	e.errs = append(e.errs, errs...)
	return e
}

func NewMultiError(err error, errs ...error) *MultiError {
	e := &MultiError{
		errs: make([]error, 0, len(errs)+1),
	}

	if err != nil {
		e.errs = append(e.errs, err)
	}
	e.errs = append(e.errs, errs...)

	return e
}

func Append(err error, errs ...error) *MultiError {
	var me *MultiError
	if errors.As(err, &me) {
		return me.add(errs...)
	}

	return NewMultiError(err).add(errs...)
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	{

		var (
			ErrNumber1 = errors.New("error 1")
			ErrNumber2 = errors.New("error 2")
			ErrNumber3 = errors.New("error 3")
		)

		var err1 error
		err1 = Append(err1, ErrNumber1)
		err1 = Append(err1, ErrNumber2)
		err1 = fmt.Errorf("internal error: %w", err1)
		err1 = Append(err1, ErrNumber3)

		if errors.Is(err1, ErrNumber1) {
			fmt.Println("found error1 with errors.Is", err1)
		}
		if errors.Is(err1, ErrNumber2) {
			fmt.Println("found error2 with errors.Is", err1)
		}
		if errors.Is(err1, ErrNumber3) {
			fmt.Println("found error3 with errors.Is", err1)
		}
		fmt.Println(err1)
	}

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
