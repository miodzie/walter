package plugin

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/mods/rss"
	"github.com/miodzie/walter/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type CommandSuite struct {
	Feed *rss.Feed
	*RssMod
	suite.Suite
	*test.SpyActions
}

func (suite *CommandSuite) SetupTest() {
	suite.RssMod = New(Context{
		Repository: rss.NewInMemRepo(),
		Parser:     &rss.StubParser{},
		Formatter:  rss.MinimalFormatter{},
	})
	suite.SpyActions = &test.SpyActions{}
	suite.RssMod.actions = suite.SpyActions
	suite.Feed = &rss.Feed{Id: 42, Name: "my_feed"}
	suite.Repository.AddFeed(suite.Feed)
}

func (suite *CommandSuite) TestSubscribeCommandNoKeywords() {
	cmd := "!subscribe my_feed"
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
		suite.Equal("", sub.Keywords)
		suite.Equal(suite.Feed.Id, sub.FeedId)
		suite.Equal("Subscribed to my_feed", suite.LastReply)
	}
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
		suite.Equal("Subscribed to my_feed with keywords: one,two. ignore: foo,baz", suite.LastReply)
	}
}

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}
