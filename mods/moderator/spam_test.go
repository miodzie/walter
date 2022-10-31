// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package moderator_test

import (
	"testing"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/moderator"
)

func TestAllCapsRegex(t *testing.T) {
	msg := seras.Message{Content: "SHUT UP DUDE"}

	if !moderator.IsSpam(msg) {
		t.Logf(`Message content: "%s" should be flagged`, msg.Content)
		t.Fail()
	}

	msg.Content = "ok then whatever"
	if moderator.IsSpam(msg) {
		t.Logf(`Message content: "%s" should not be flagged`, msg.Content)
		t.Fail()
	}
}
