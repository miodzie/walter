package rss

import (
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"strings"
	"testing"
)

type AnnouncementOrganizerSuite struct {
	organizer *AnnouncementOrganizer
}

func (s *AnnouncementOrganizerSuite) PreTest(t *td.T, name string) error {
	s.organizer = &AnnouncementOrganizer{}
	return nil
}

func (s *AnnouncementOrganizerSuite) TestOrganizeGroupsByGUIDAndChannel(assert *td.T) {
	organizer := AnnouncementOrganizer{}
	item := Item{Title: "foo", Link: "blog.golang.org"}
	bob := Notification{User: "bob", Channel: "#general", Item: item}
	carl := Notification{User: "carl", Channel: "#general", Item: item}

	announces := organizer.Organize([]Notification{bob, carl})

	if assert.Len(announces, 1) {
		assert.Cmp(announces[0].Room, "#general")
		assert.Cmp(announces[0].Message, "foo - blog.golang.org : bob,carl")
	}
}

func (s *AnnouncementOrganizerSuite) TestOrganizeSameChannelDifferentGUID(assert *td.T) {
	organizer := AnnouncementOrganizer{}
	bob := Notification{User: "bob", Channel: "#general",
		Item: Item{
			GUID:  "A",
			Title: "hi bob",
		}}
	carl := Notification{User: "carl", Channel: "#general",
		Item: Item{
			GUID:  "B",
			Title: "hi carl",
		}}

	announces := organizer.Organize([]Notification{bob, carl})

	if assert.Len(announces, 2) {
		var b, c Announcement
		// Oh, the pains of hashmaps.
		if announces[0].users[0] == "bob" {
			b, c = announces[0], announces[1]
		} else {
			b, c = announces[1], announces[0]
		}
		assert.True(strings.Contains(b.Message, "hi bob"))
		assert.True(strings.Contains(c.Message, "hi carl"))
	}
}

func TestAnnouncementOrganizerSuite(t *testing.T) {
	tdsuite.Run(t, new(AnnouncementOrganizerSuite))
}
