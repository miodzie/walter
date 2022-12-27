package stocks

import (
	"fmt"
)

type UnknownSecurity struct {
	Symbol string
}

func (u UnknownSecurity) Error() string {
	return fmt.Sprintf("unknown security: %s", u.Symbol)
}

type Ticker interface {
	PriceFor(tickerSymbol string) (float64, error)
}
