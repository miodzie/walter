package rss

import (
	"fmt"
	"strings"
)

// EXCUSE ME! I HAVE TO MAKE AN ANNOUNCEMENT!

// By having the Notification be separate,
// I can keep the option of delivering notifications instead of announcements through
// PMs, or possible digest emails/links.

type Announcer interface {
	Announce([]Announcement) error
}

type Announcement struct {
	Message    string
	Room       string
	OnDelivery func() error

	users []string
}

// Fetch Feeds  -> Create Notification -> Organize into Announcements
// -> Announce (Decorate Announcer for filters)

type AnnouncementOrganizer struct{}

func (o *AnnouncementOrganizer) Organize(notes []Notification) []Announcement {
	var announces []Announcement
	seen := make(map[string]*Announcement)
	for _, n := range notes {
		a, exists := seen[o.key(n)]
		if exists {
			a.users = append(a.users, n.User)
			a.Message = o.formatMsg(n, a.users)
		} else {
			a2 := &Announcement{
				Room:  n.Channel,
				users: []string{n.User},
			}
			a2.Message = o.formatMsg(n, a2.users)
			seen[o.key(n)] = a2
		}
	}
	for _, a := range seen {
		announces = append(announces, *a)
	}

	return announces
}
func (o *AnnouncementOrganizer) key(n Notification) string {
	return n.Channel + n.Item.GUID
}

// TODO: Consider replacing with formatter.
func (o *AnnouncementOrganizer) formatMsg(n Notification, users []string) string {
	return fmt.Sprintf(
		"%s - %s : %s",
		n.Item.Title, n.Item.Link, strings.Join(users, ","))
}
