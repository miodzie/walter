package rss

// Matcher matches Subscriptions to Feed Items.
// Returns Notifications that satisfies the Subscription requirements.
type Matcher struct {
	subs []Subscription
}

func NewMatcher(subs []Subscription) *Matcher {
	return &Matcher{subs: subs}
}

func (m *Matcher) Match(items []Item) (matches []Notification) {
	for _, sub := range m.subs {
		matches = append(matches, findMatches(sub, items)...)
	}
	return
}

func findMatches(sub Subscription, items []Item) (matches []Notification) {
	for _, item := range items {
		if !sub.ShouldIgnore(item) {
			note := Notification{
				Channel: sub.Channel,
				Feed:    *sub.Feed,
				Item:    item,
				User:    sub.User,
			}
			matches = append(matches, note)
		}
	}
	return matches
}
