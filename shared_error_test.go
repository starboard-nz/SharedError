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

func TestSharedError(t *testing.T) {
	sharedErr := sharederror.SharedError{}
	sharedErr.Store(fmt.Errorf("some error"))

	if !sharedErr.Triggered() {
		t.Error("expecting shared-error to trigger on stored error")
	}
}
