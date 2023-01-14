package rss

import (
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

func TestMatcherMatchesNotificationsToSubscriptions(t *testing.T) {
	var subs = []Subscription{{User: "bob", Feed: &UserFeed{}}}
	matcher := NewMatcher(subs)
	posting := Item{GUID: "1234", Title: "New Blog Notification"}

	mail := matcher.Match([]Item{posting})

	expected := Notification{
		User:    "bob",
		Feed:    UserFeed{},
		Item:    posting,
		Channel: "",
	}
	if td.Cmp(t, mail, td.Len(1)) {
		td.Cmp(t, mail[0], expected)
	}
}
