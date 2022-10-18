/**
 * Copyright (c) 2021, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE.TXT in the root directory of this source tree.
 */

package sharederror

import (
	"errors"
	"fmt"
	"sync"
)

// SharedError implements goroutine-safe error handling.
// Multiple concurrent goroutines can share a SharedError variable to reports errors
// to the calling function.
type SharedError struct {
	err  []error
	lock sync.Mutex
}

// NewSharedError creates new shared-error.
func NewSharedError() *SharedError {
	return &SharedError{}
}

// Error implements the error interface, you can return a SharedError as an error.
func (s *SharedError) Error() string {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.err) == 0 {
		return ""
	}

	if len(s.err) == 1 {
		return s.err[0].Error()
	}

	errorStr := ""

	for i, err := range s.err {
		if i != 0 {
			errorStr += " / "
		}

		errorStr += fmt.Sprintf("error %d: %v", i, err)
	}

	return errorStr
}

// Triggered returns true if an error is stored in the ShareError, or false for no error.
func (s *SharedError) Triggered() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.err) == 0 {
		return false
	}

	return true
}

// Store stores an error condition in SharedError.
func (s *SharedError) Store(err error) {
	if err == nil {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.err = append(s.err, err)
}

// Store stores an error condition in SharedError.
func (s *SharedError) Storef(format string, args ...interface{}) {
	if format == "" {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.err = append(s.err, fmt.Errorf(format, args...))
}

// Err returns SharedError as an error when triggered.
func (s *SharedError) Err() error {
	if s.Triggered() {
		return s
	} else {
		return nil
	}
}

// Errors returns SharedError errors if any.
func (s *SharedError) Errors() []error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.err
}

// Has checks if any error of SharedError contains target error through errors.Is.
func (s *SharedError) Has(target error) bool {
	for _, err := range s.Errors() {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// HasOnly checks if all errors of SharedError contains target error through errors.Is.
func (s *SharedError) HasOnly(target error) bool {
	for _, err := range s.Errors() {
		if !errors.Is(err, target) {
			return false
		}
	}
	return true
}
