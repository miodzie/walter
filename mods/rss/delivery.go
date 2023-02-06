package rss

func ThrottleByChannel(notes chan Notification, maxPerChannel int) chan Notification {
	wrapped := make(chan Notification)
	go func() {
		seen := make(map[string]int)
		for n := range notes {
			seen[n.Channel]++
			if seen[n.Channel] != maxPerChannel {
				wrapped <- n
			}
		}
		close(wrapped)
	}()
	return wrapped
}
