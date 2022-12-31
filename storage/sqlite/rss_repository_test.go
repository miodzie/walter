package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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

func TestLongRssRepositorySuite(t *testing.T) {
	suite.Run(t, new(RssRepositorySuite))
}
