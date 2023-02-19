package delivery

// Digest is a Deliverable that pools given deliveries until a configured threshold,
// then delivers them all at once.
// Plans are for a digest by Email, and/or by viewable by a private link.

// Digest should probably accept chan Notification, not chan Deliverable.
// This allows more nuance into formatting, not sure yet.

type Digest struct {
}

func (d Digest) Address() string {
	//TODO implement me
	panic("implement me")
}

func (d Digest) Deliver(deliver func(address string, content string) error) {
	//TODO implement me
	panic("implement me")
}

/*

	- Subscription/Feed Management
	- (Processing?),
	- Delivery


*/
