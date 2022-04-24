package policing_test

import (
	"testing"

	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/policing"
)

func TestAllCapsRegex(t *testing.T) {
  msg := seras.Message{Content: "SHUT UP DUDE"}

  if (!policing.IsSpam(msg)) {
    t.Logf(`Message content: "%s" should be flagged`, msg.Content)
    t.Fail()
  }

  msg.Content = "ok then whatever"
  if (policing.IsSpam(msg)) {
    t.Logf(`Message content: "%s" should not be flagged`, msg.Content)
    t.Fail()
  }
}
