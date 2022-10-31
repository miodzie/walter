// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package trigger

import "testing"

func TestTrigger_Check_contains_word(t *testing.T) {
	strs := []string{
		"I own a cow.",
		"I OwN a cow.",
	}

	trig := &Trigger{Word: "own"}

	for _, s := range strs {
		if !trig.Check(s) {
			t.Fail()
		}
	}
}

func TestTrigger_Check_doesnt_contain_word(t *testing.T) {
	str := "I own a cow."

	trig := &Trigger{Word: "chicken"}

	if trig.Check(str) {
		t.Fail()
	}
}
