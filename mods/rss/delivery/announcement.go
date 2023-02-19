package delivery

import (
	"fmt"
	"strings"
)

// EXCUSE ME! I HAVE TO MAKE AN ANNOUNCEMENT!

// By having the Notification be separate,
// I can keep the option of delivering notifications instead of announcements through
// PMs, or possible digest emails/links.

// I can abstract this out into a pipeline that returns a channel of
// Notifications, this enables more modularity.
// I can then have a RealTimeProcessor,
//that's constantly polling and sending new Notifications fresh off the channel.
// While this aggregate into announcements can be a separate pipeline,
//off the same base.

type Announcement struct {
	message      string
	room         string
	deliveryHook func() error

	users []string
}

func (a Announcement) Address() string {
	return a.room
}

func (a Announcement) Deliver(deliver func(address string, content string) error) {
	if deliver(a.room, a.message) != nil {
		_ = a.deliveryHook()
	}
}

// Fetch Feeds  -> Create Notification -> Organize into Announcements
// -> Announce (Decorate Announcer for filters)

// AnnouncementOrganizer organizes new Subscription Notifications into grouped
// Announcements, by Channel(room) and the Item's GUID.
// TODO: Add Formatter to creation?
type AnnouncementOrganizer struct{}

func (o *AnnouncementOrganizer) Organize(notes []Notification) []Announcement {
	var announces []Announcement
	seen := make(map[string]*Announcement)
	for _, n := range notes {
		a, exists := seen[o.key(n)]
		if exists {
			a.users = append(a.users, n.User)
			a.message = o.formatMsg(n, a.users)
		} else {
			a2 := &Announcement{
				room:  n.Channel,
				users: []string{n.User},
			}
			a2.message = o.formatMsg(n, a2.users)
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

func (o *AnnouncementOrganizer) formatMsg(n Notification, users []string) string {
	return fmt.Sprintf(
		"%s - %s : %s",
		n.Item.Title, n.Item.Link, strings.Join(users, ","))
}
