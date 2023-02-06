package rss

import (
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

func TestThrottleByChannel(t *testing.T) {
	s := &ChannelThrottler{Max: 1}
	notes := make(chan Notification, 2)
	notes <- Notification{User: "jacob", Channel: "#go"}
	notes <- Notification{User: "abraham", Channel: "#go"}
	close(notes)

	notes = s.Throttle(notes)

	td.Cmp(t, (<-notes).User, "abraham")
	td.Cmp(t, (<-notes).User, "")
}
