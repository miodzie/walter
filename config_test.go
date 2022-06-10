package seras

import (
	"testing"
)

func TestParseToml(t *testing.T) {
	cfg, err := ParseToml("config.toml")
	if err != nil {
		t.Error(err)
	}

	if len(cfg.Bots) != 2 {
		t.Fail()
	}
	d, ok := cfg.Bots["discord"]
	if !ok {
		t.Fail()
	}
	if d["type"] != "discord" {
		t.Fail()
	}
}
