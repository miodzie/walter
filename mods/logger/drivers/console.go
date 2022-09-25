package drivers

import (
	"fmt"
	"github.com/miodzie/seras"
)

type ConsoleLogger struct {
}

func (l ConsoleLogger) Log(m seras.Message) error {
	fmt.Printf("[%s] <%s>: %s\n", m.ConnectionName, m.Author.Nick, m.Content)
	return nil
}
