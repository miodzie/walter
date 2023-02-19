package delivery

import (
	"fmt"
	"github.com/miodzie/walter/mods/rss"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotification_String(t *testing.T) {
	item := rss.Item{Title: "New Cool Blog Notification",
		Description: "Lorem ispsum but cooler.",
		Link:        "http://localhost"}
	notification := Notification{Item: item, User: "Abraham"}
	expected := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		item.Title,
		item.DescriptionTruncated(),
		item.Link,
		notification.User)

	assert.Equal(t, expected, notification.String())
	//fmt.Println(expected)
}

func TestMinimalFormatter_Format(t *testing.T) {
	formatter := MinimalFormatter{}
	item := rss.Item{Title: "New Cool Blog Notification",
		Description: "Lorem ispsum but cooler.",
		Link:        "http://localhost"}
	notification := Notification{Item: item, User: "Abraham"}

	expected := fmt.Sprintf(
		"%s - %s : %s",
		item.Title, item.Link, notification.User,
	)

	assert.Equal(t, expected, formatter.Format(notification))
	//fmt.Println(expected)
}
