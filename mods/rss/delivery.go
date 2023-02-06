package rss

type ChannelThrottler struct {
	Max int
}

func (t *ChannelThrottler) Throttle(notes chan Notification) chan Notification {
	wrapped := make(chan Notification)
	go func() {
		seen := make(map[string]int)
		for n := range notes {
			seen[n.Channel]++
			if seen[n.Channel] != t.Max {
				wrapped <- n
			}
		}
		close(wrapped)
	}()
	return wrapped
}
