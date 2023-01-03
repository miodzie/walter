package plugin

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"walter"
	"walter/mods/rss/internal/usecases"
	"walter/test"
)

type CommandSuite struct {
	*RssMod
	*test.SpyActions
	suite.Suite
}

func (suite *CommandSuite) SetupTest() {
	suite.RssMod = New(Services{
		Repository: usecases.NewInMemRepo(),
		Parser:     &usecases.StubParser{},
		Formatter:  usecases.MinimalFormatter{},
	})
	suite.SpyActions = &test.SpyActions{}
	suite.RssMod.actions = suite.SpyActions
	newFeedRequest := usecases.AddFeedRequest{Name: "my_feed"}
	_, err := usecases.NewAddFeed(suite.Repository).Add(newFeedRequest)
	suite.NoError(err)
}

func (suite *CommandSuite) TestRemoveFeed() {
	cmd := "!remove_feed my_feed"
	msg := walter.Message{
		Content: cmd, Arguments: strings.Split(cmd, " "),
		Target: "##feeds",
		Author: walter.Author{Id: "admin", Mention: "author_mention"},
	}
	suite.AdminUserId = msg.Author.Id

	suite.removeFeed(msg)

	_, err := suite.Repository.FeedByName("my_feed")
	suite.ErrorIs(err, usecases.FeedNotFoundError)
}

func (suite *CommandSuite) TestSubscribeCommandNoKeywords() {
	cmd := "!subscribe my_feed"
	msg := walter.Message{
		Content: cmd, Arguments: strings.Split(cmd, " "),
		Target: "##feeds",
		Author: walter.Author{Id: "author_id", Mention: "author_mention"},
	}

	suite.subscribe(msg)

	subs, err := suite.Repository.Subs(usecases.SearchParams{})
	assert.Nil(suite.T(), err)
	if assert.Len(suite.T(), subs, 1) {
		sub := subs[0]
		suite.Equal("author_mention", sub.User)
		suite.Equal("##feeds", sub.Channel)
		suite.Equal("", sub.Keywords)
		//suite.Equal(suite.Feed.Id, sub.FeedId)
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

	subs, err := suite.Repository.Subs(usecases.SearchParams{})
	assert.Nil(suite.T(), err)
	if assert.Len(suite.T(), subs, 1) {
		sub := subs[0]
		suite.Equal("author_mention", sub.User)
		suite.Equal("##feeds", sub.Channel)
		suite.Equal("foo,bar", sub.Keywords)
		//suite.Equal(suite.Feed.Id, sub.FeedId)
	}
}

func (suite *CommandSuite) TestSubscribeCommandParsesIgnoreWords() {
	cmd := "!subscribe my_feed -keywords=one,two -ignore=foo,baz"
	msg := walter.Message{Content: cmd, Arguments: strings.Split(cmd, " ")}

	suite.subscribe(msg)

	subs, err := suite.Repository.Subs(usecases.SearchParams{})
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
