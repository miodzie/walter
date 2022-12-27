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

func (suite *CommandSuite) SetupTest() {
	suite.RssMod = New(Context{
		Repository: rss.NewInMemRepo(),
		Parser:     &rss.StubParser{},
		Formatter:  rss.MinimalFormatter{},
	})
	suite.SpyActions = &SpyActions{}
	suite.RssMod.actions = suite.SpyActions
	suite.Feed = &rss.Feed{Id: 42, Name: "my_feed"}
	suite.Repository.AddFeed(suite.Feed)
}

func (suite *CommandSuite) TestSubscribeCommandSubscribesWithKeywords() {
	cmd := "!subscribe my_feed -keywords=foo,bar"
	msg := walter.Message{
		Content: cmd, Arguments: strings.Split(cmd, " "),
		Target: "##feeds",
		Author: walter.Author{Id: "author_id", Mention: "author_mention"},
	}

	suite.subscribe(msg)

	subs, err := suite.Repository.Subs(rss.SearchParams{})
	assert.Nil(suite.T(), err)
	if assert.Len(suite.T(), subs, 1) {
		sub := subs[0]
		suite.Equal("author_mention", sub.User)
		suite.Equal("##feeds", sub.Channel)
		suite.Equal("foo,bar", sub.Keywords)
		suite.Equal(suite.Feed.Id, sub.FeedId)
	}
}

func (suite *CommandSuite) TestSubscribeCommandParsesIgnoreWords() {
	cmd := "!subscribe my_feed -keywords=one,two -ignore=foo,baz"
	msg := walter.Message{Content: cmd, Arguments: strings.Split(cmd, " ")}

	suite.subscribe(msg)

	subs, err := suite.Repository.Subs(rss.SearchParams{})
	suite.Nil(err)
	if suite.Len(subs, 1) {
		sub := subs[0]
		suite.Equal("one,two", sub.Keywords)
		suite.Equal("foo,baz", sub.Ignore)
		suite.Equal("Subscribed to my_feed with keywords: one,two. ignore: foo,baz", suite.SpyActions.ReplyMsg)
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
