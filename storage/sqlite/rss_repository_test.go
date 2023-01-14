package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/miodzie/walter/mods/rss"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RssRepositorySuite struct {
	repository *RssRepository
	db         *sql.DB
	suite.Suite
}

func (test *RssRepositorySuite) SetupTest() {
	var err error
	test.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	if err = migrate(test.db); err != nil {
		panic(err)
	}
	test.repository = NewRssRepository(test.db)
}

func (test *RssRepositorySuite) TestFeeds() {
	_, err := test.repository.Feeds()
	fmt.Println(err)
	test.Nil(err)
}

func (test *RssRepositorySuite) SkipTestRemoveFeed() {
	feed := &rss.UserFeed{Name: "foo"}
	test.NoError(test.repository.AddFeed(feed))

	err := test.repository.RemoveFeed("foo")

	if test.NoError(err) {
		f, err := test.repository.Feeds()
		test.NoError(err)
		test.Empty(f)
	}
}

func TestLongRssRepositorySuite(t *testing.T) {
	suite.Run(t, new(RssRepositorySuite))
}
