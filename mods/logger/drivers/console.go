package drivers

import (
	"fmt"
	"github.com/miodzie/seras"
)

type ConsoleLogger struct {
}

func (l ConsoleLogger) Log(m seras.Message) error {
	if m.Content != "" {
		if m.Target == "" {
			fmt.Printf("[%s] <%s>: %s\n", m.ConnectionName, m.Author.Nick, m.Content)
			return nil
		}
		fmt.Printf("[%s] (%s) <%s>: %s\n", m.ConnectionName, m.Target, m.Author.Nick, m.Content)
		return nil
	}

	fmt.Printf("[%s]: %s\n", m.ConnectionName, m.Raw)
	return nil
}
