package rss

import (
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

type AnnouncementOrganizerSuite struct {
	organizer *AnnouncementOrganizer
}

func (s *AnnouncementOrganizerSuite) PreTest(t *td.T, name string) error {
	s.organizer = &AnnouncementOrganizer{}
	return nil
}

func (s *AnnouncementOrganizerSuite) TestOrganizeReturnsAnnouncements(assert *td.T) {
	organizer := AnnouncementOrganizer{}
	bob := Notification{
		User:    "bob",
		Channel: "#general",
	}

	announces := organizer.Organize([]Notification{bob})

	if assert.Len(announces, 1) {
		assert.Cmp(announces[0].Room, "#general")
		assert.Cmp(announces[0].Message, bob.String())
	}
}

func TestAnnouncementOrganizerSuite(t *testing.T) {
	tdsuite.Run(t, new(AnnouncementOrganizerSuite))
}
