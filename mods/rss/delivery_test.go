package rss

import (
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

func TestThrottleByChannel(t *testing.T) {
	notes := make(chan Deliverable, 3)
	notes <- Notification{User: "jacob", Channel: "#go"}
	notes <- Notification{User: "issac", Channel: "#go"}
	notes <- Notification{User: "abraham", Channel: "#go"}

	notes = ThrottleByChannel(notes, 2)

	td.Cmp(t, (<-notes).Address(), "#go")
	td.Cmp(t, (<-notes).Address(), "#go")
	td.Cmp(t, <-notes, td.Empty())
}
