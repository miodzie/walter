package usecases

import (
	"fmt"
	"github.com/miodzie/walter/mods/rss/internal/internal/domain"
	"strings"
)

type Formatter interface {
	Format(domain.Notification) string
}

type DefaultFormatter struct {
}

func (d DefaultFormatter) Format(notification domain.Notification) string {
	return notification.String()
}

type MinimalFormatter struct {
}

func (m MinimalFormatter) Format(n domain.Notification) string {
	i := n.Item
	return fmt.Sprintf(
		"%s - %s : %s",
		i.Title, i.Link, strings.Join(n.Users, ","),
	)
}
