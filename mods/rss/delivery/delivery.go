package delivery

func ThrottleByChannel(notes chan Deliverable, maxPerChannel int) chan Deliverable {
	wrapped := make(chan Deliverable)
	go func() {
		seen := make(map[string]int)
		for n := range notes {
			if seen[n.Address()] >= maxPerChannel {
				close(wrapped)
				return
			}
			wrapped <- n
			seen[n.Address()]++
		}
		close(wrapped)
	}()
	return wrapped
}
