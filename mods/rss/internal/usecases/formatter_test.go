package usecases

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"walter/mods/rss/internal/internal/domain"
)

func TestMinimalFormatter_Format(t *testing.T) {
	formatter := MinimalFormatter{}
	item := domain.ParsedItem{Title: "New Cool Blog Post",
		Description: "Lorem ispsum but cooler.",
		Link:        "http://localhost"}
	notification := domain.Notification{Item: item, Users: []string{"Abraham", "Isaac", "Jacob"}}

	expected := fmt.Sprintf(
		"%s - %s : %s",
		item.Title, item.Link, strings.Join(notification.Users, ","))

	assert.Equal(t, expected, formatter.Format(notification))
	//fmt.Println(expected)
}
