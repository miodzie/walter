package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNotification_String(t *testing.T) {
	item := Item{Title: "New Cool Blog Post",
		Description: "Lorem ispsum but cooler.",
		Link:        "http://localhost"}
	notification := Notification{Item: item, Users: []string{"Abraham", "Isaac", "Jacob"}}
	expected := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		item.Title,
		item.DescriptionTruncated(),
		item.Link,
		strings.Join(notification.Users, ", "))

	assert.Equal(t, expected, notification.String())
	//fmt.Println(expected)
}
