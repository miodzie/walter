package test

import (
	"github.com/miodzie/walter"
	"time"
)

// SpyActions provides some helpful commands when testing expected Action calls on a module.
type SpyActions struct {
	LastReply   string
	AdminUserId string
}

func (s *SpyActions) Send(message walter.Message) error {
	return nil
}

func (s *SpyActions) Reply(message walter.Message, reply string) error {
	s.LastReply = reply
	return nil
}

func (s *SpyActions) Bold(s2 string) string {
	return ""
}

func (s *SpyActions) Italicize(s2 string) string {
	return ""
}

func (s *SpyActions) IsAdmin(userId string) bool {
	return userId == s.AdminUserId
}

func (s *SpyActions) TimeoutUser(channel string, user string, until time.Time) error {
	return nil
}
