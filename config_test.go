package seras

import (
	"fmt"
	"testing"
)

func TestParseToml(t *testing.T) {
  cfg, err := ParseToml("config.toml")
  if err != nil {
    t.Error(err)
  }

  t.Fail()
  fmt.Println(cfg)
}
