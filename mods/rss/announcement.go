package rss

// EXCUSE ME! I HAVE TO MAKE AN ANNOUNCEMENT!

// Fetch Feeds -> Notification -> Announcements -> Filters -> Deliver
// Filters:
// - Channel Limit

// By having the Notification be separate, I can later port this to self-hosted stuff.

type Announcement struct {
	Message    string
	Room       string
	OnDelivery func() error
}

// Fetch Feeds  -> Create Notification -> Organize into Announcements -> Deliver

// []Notifications -> Announcement

type AnnouncementOrganizer struct{}

func (o *AnnouncementOrganizer) Organize(notes []Notification) []Announcement {
	var announces []Announcement

	for _, n := range notes {
		a := Announcement{Room: n.Channel, Message: n.String()}
		announces = append(announces, a)
	}

	return announces
}
