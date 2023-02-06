package rss

func ThrottleByChannel(notes chan Deliverable, maxPerChannel int) chan Deliverable {
	wrapped := make(chan Deliverable)
	go func() {
		seen := make(map[string]int)
		for n := range notes {
			seen[n.Address()]++
			if seen[n.Address()] != maxPerChannel {
				wrapped <- n
			}
		}
		close(wrapped)
	}()
	return wrapped
}
