package sqlite

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Setup("testing.db")
	code := m.Run()
	err := os.Remove("testing.db")
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
