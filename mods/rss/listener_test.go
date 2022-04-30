package rss_test

import (
	"testing"

	"github.com/miodzie/seras/storage/memory"
)

var repo memory.InMemListenerRepository

func TestThing(t *testing.T) {
    listeners := repo.All()
    if len(listeners) == 0 {
        t.Fail()
    }
}
