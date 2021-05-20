/**
 * Copyright (c) 2021, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE.TXT in the root directory of this source tree.
 */

package sharederror

import (
	"fmt"
	"sync"
)

// SharedError implements goroutine-safe error handling.
// Multiple concurrent goroutines can share a SharedError variable to reports errors
// to the calling function.
type SharedError struct {
	err []error
	lock sync.Mutex
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
