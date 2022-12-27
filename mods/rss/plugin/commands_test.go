package plugin

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/rss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type CommandSuite struct {
	Feed *rss.Feed
	*RssMod
	suite.Suite
	SpyActions *SpyActions
}

func (s *CommandSuite) SetupTest() {
	s.RssMod = New(Context{
		Repository: rss.NewInMemRepo(),
		Parser:     &rss.StubParser{},
		Formatter:  rss.MinimalFormatter{},
	})
	s.SpyActions = &SpyActions{}
	s.RssMod.actions = s.SpyActions
	s.Feed = &rss.Feed{Id: 42, Name: "my_feed"}
	s.Repository.AddFeed(s.Feed)
}

func (s *CommandSuite) TestSubscribeCommandSubscribesWithKeywords() {
	cmd := "!subscribe my_feed -keywords=foo,bar"
	msg := walter.Message{
		Content:   cmd,
		Arguments: strings.Split(cmd, " "),
		Author:    walter.Author{Id: "author_id", Mention: "author_mention"},
		Target:    "##feeds",
	}

	s.subscribe(msg)

	subs, err := s.Repository.Subs(rss.SearchParams{})
	assert.Nil(s.T(), err)
	if assert.Len(s.T(), subs, 1) {
		sub := subs[0]
		assert.Equal(s.T(), "author_mention", sub.User)
		assert.Equal(s.T(), "##feeds", sub.Channel)
		assert.Equal(s.T(), "foo,bar", sub.Keywords)
		assert.Equal(s.T(), s.Feed.Id, sub.FeedId)
	}
}

func (s *CommandSuite) TestSubscribeCommandParsesIgnoreWords() {
	cmd := "!subscribe my_feed -keywords=one,two -ignore=foo,baz"
	msg := walter.Message{Content: cmd, Arguments: strings.Split(cmd, " ")}

	s.subscribe(msg)

	subs, err := s.Repository.Subs(rss.SearchParams{})
	assert.Nil(s.T(), err)
	if assert.Len(s.T(), subs, 1) {
		sub := subs[0]
		assert.Equal(s.T(), "one,two", sub.Keywords)
		assert.Equal(s.T(), "foo,baz", sub.Ignore)
		assert.Equal(s.T(), "Subscribed to my_feed with keywords: one,two. ignore: foo,baz", s.SpyActions.ReplyMsg)
	}
}

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

// //////////////////////////////////////////////////////////////////////////
type SpyActions struct {
	ReplyMsg string
}

func (s SpyActions) Send(message walter.Message) error {
	return nil
}

func (s *SpyActions) Reply(message walter.Message, reply string) error {
	s.ReplyMsg = reply
	return nil
}

func (s SpyActions) Bold(s2 string) string {
	return ""
}

func (s SpyActions) Italicize(s2 string) string {
	return ""
}

func (s SpyActions) IsAdmin(userId string) bool {
	return false
}

func (s SpyActions) TimeoutUser(channel string, user string, until time.Time) error {
	return nil
}
