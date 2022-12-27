// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package art

import (
	"fmt"
	"testing"
)

func SkipTestArt_gm(t *testing.T) {
	art := &Picture{Art: gm}
	fmt.Println("test")
	for !art.Completed() {
		t.Log(art.NextLine())
	}
	//t.Fail()
}
