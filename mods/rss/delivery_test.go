package rss

import (
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

func TestThrottleByChannel(t *testing.T) {
	notes := make(chan Notification, 2)
	notes <- Notification{User: "jacob", Channel: "#go"}
	notes <- Notification{User: "abraham", Channel: "#go"}
	close(notes)

	notes = ThrottleByChannel(notes, 1)

	td.Cmp(t, (<-notes).User, "abraham")
	td.Cmp(t, (<-notes).User, "")
}