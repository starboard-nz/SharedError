/**
 * Copyright (c) 2022, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE.TXT in the root directory of this source tree.
 */

package sharederror_test

import (
	"fmt"
	"testing"

	"github.com/starboard-nz/sharederror"
)

func TestNewSharedError(t *testing.T) {
	var sharedErr interface{} = sharederror.NewSharedError()
	_, ok := (sharedErr).(error)
	if !ok {
		t.Errorf("expected shared error to implement error")
	}
}

func TestSharedErrorTriggered(t *testing.T) {
	var sharedErr = sharederror.NewSharedError()

	t.Log("initial")
	if sharedErr.Triggered() {
		t.Errorf("expected shared error not to be triggered")
	}

	t.Log("store error")
	sharedErr.Store(fmt.Errorf("some error"))
	if !sharedErr.Triggered() {
		t.Errorf("expected shared error to be triggered")
	}
}

func TestSharedErrorErr(t *testing.T) {
	var sharedErr = sharederror.NewSharedError()

	t.Log("initial")
	if err := sharedErr.Err(); err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	t.Log("store error")
	sharedErr.Store(fmt.Errorf("some error"))
	if err := sharedErr.Err(); err == nil {
		t.Errorf("expected error after store, got nil")
	}
}

func TestSharedErrorErrors(t *testing.T) {
	var sharedErr = sharederror.NewSharedError()

	t.Log("initial")
	if errs := sharedErr.Errors(); len(errs) > 0 {
		t.Errorf("expected no errors, got %d", len(errs))
	}

	n := 30
	t.Log("store errors")
	for i := 0; i < n; i++ {
		sharedErr.Store(fmt.Errorf("some error"))
	}

	t.Log("after errors")
	if errs := sharedErr.Errors(); len(errs) != n {
		t.Errorf("expected %d stored errors, got %d", n, len(errs))
	}
}

type targetError struct{}

func (e *targetError) Error() string {
	return "some error"
}

func TestSharedErrorIsAny(t *testing.T) {
	var sharedErr = sharederror.NewSharedError()

	t.Log("initial")
	if sharedErr.IsAny(&targetError{}) {
		t.Errorf("expected no match on target error, got match")
	}

	t.Log("store other error")
	sharedErr.Store(fmt.Errorf("some other error"))
	if sharedErr.IsAny(&targetError{}) {
		t.Errorf("expected no match on target error, got match")
	}

	t.Log("store target error")
	sharedErr.Store(&targetError{})
	if !sharedErr.IsAny(&targetError{}) {
		t.Errorf("expected match on target error, got no matches")
	}
}

func TestSharedErrorIsAll(t *testing.T) {
	var sharedErr = sharederror.NewSharedError()

	t.Log("initial")
	if sharedErr.IsAll(&targetError{}) {
		t.Errorf("expected no match on target error, got match")
	}

	n := 30
	t.Log("store target errors")
	for i := 0; i < n; i++ {
		sharedErr.Store(&targetError{})
	}

	t.Log("is all target errors")
	if !sharedErr.IsAll(&targetError{}) {
		t.Errorf("expected all errors match on target error, got no matches")
	}

	t.Log("store other error")
	sharedErr.Store(fmt.Errorf("some other error"))
	if sharedErr.IsAll(&targetError{}) {
		t.Errorf("expected not all errors match on target error, got all matches")
	}
}

func TestReset(t *testing.T) {
	var sharedErr = sharederror.NewSharedError()

	err := sharedErr.Reset()
	if err != nil {
		t.Errorf("incorrect result")
	}

	t.Log("store error")
	sharedErr.Store(fmt.Errorf("some error"))

	err = sharedErr.Reset()
	if err == nil {
		t.Errorf("incorrect result")
	}

	err = sharedErr.Reset()
	if err != nil {
		t.Errorf("incorrect result")
	}
}
