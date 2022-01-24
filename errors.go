package gobl

import (
	"errors"
	"fmt"
)

// Error provides a structure to better be able to make error comparisons.
// The contents can also be serialised as JSON ready to send to a client
// if needed.
type Error struct {
	Code  string `json:"code"`
	Cause error  `json:"cause"`
}

var (
	// ErrNoRegion is used when the envelope is missing a region.
	ErrNoRegion = NewError("no-region")

	// ErrNoDocument is provided when the envelope does not contain a
	// document payload.
	ErrNoDocument = NewError("no-document")

	// ErrValidation is used when a document fails a validation request.
	ErrValidation = NewError("validation")

	// ErrCalculation wraps around errors that we're generated during a
	// call to perform calculations on a document.
	ErrCalculation = NewError("calculation")

	// ErrMarshal is provided when there has been a problem attempting to encode
	// or marshal an object, usually into JSON.
	ErrMarshal = NewError("marshal")

	// ErrSignature identifies an issue related to signatures.
	ErrSignature = NewError("signature")

	// ErrInternal is a "catch-all" for errors that are not expected.
	ErrInternal = NewError("internal")
)

// NewError provides a new error with a code that is meant to provide
// a context.
func NewError(code string) *Error {
	return &Error{Code: code}
}

// Error provides a string representation of the error.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Cause.Error())
	}
	return e.Code
}

// WithCause is used to copy and add an underlying error to this one.
func (e *Error) WithCause(err error) *Error {
	ne := e.copy()
	ne.Cause = err
	return ne
}

// WithErrorf wraps around the `fmt.Errorf` call to provide a more meaningful
// error in the context.
func (e *Error) WithErrorf(format string, a ...interface{}) *Error {
	ne := e.copy()
	ne.Cause = fmt.Errorf(format, a...)
	return ne
}

func (e *Error) copy() *Error {
	ne := new(Error)
	*ne = *e
	return ne
}

// Is checks to see if the target error matches the current error or
// part of the chain.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.Cause, target)
	}
	return e.Code == t.Code
}
