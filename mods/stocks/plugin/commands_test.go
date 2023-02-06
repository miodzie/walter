package plugin

import (
	"github.com/miodzie/walter"
	"github.com/miodzie/walter/test"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type CommandSuite struct {
	*StockMod
	*test.SpyActions
	suite.Suite
}

func (s *CommandSuite) SetupTest() {
	s.StockMod = &StockMod{Ticker: &StubTicker{}}
	s.SpyActions = &test.SpyActions{}
	s.StockMod.Actions = s.SpyActions
}

//func (s *CommandSuite) TestPriceUnknownSecurity() {
//	cmd := "!price POTATO"
//	msg := walter.message{Content: cmd, Arguments: strings.Split(cmd, " ")}
//
//	s.getPrice(msg)
//
//	s.Equal("Unknown security: 'POTATO'", s.LastReply)
//}

func (s *CommandSuite) TestPriceCommand() {
	cmd := "!price STONK"
	msg := walter.Message{Content: cmd, Arguments: strings.Split(cmd, " ")}

	s.getPrice(msg)

	s.Equal("$STONK: 10.00", s.LastReply)
}

func (s *CommandSuite) TestPriceEmpty() {
	msg := walter.Message{Content: "", Arguments: strings.Split("", " ")}
	s.getPrice(msg)
	s.Equal("", s.LastReply)
}

func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

// //////////////////////////////////////////////////////////////////////////

type StubTicker struct {
}

func (s StubTicker) PriceFor(symbol string) (float64, error) {
	return 10, nil
}
