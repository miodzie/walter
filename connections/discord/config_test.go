package discord

import "testing"

func TestParseConfig_parses_into_Config_struct(t *testing.T) {
	val := make(map[string]interface{})
	val["type"] = "discord"
	val["token"] = "my_token"
	val["admins"] = []interface{}{"alice"}
	val["mods"] = []interface{}{"myplugin"}

	cfg, err := ParseConfig(val)
	if err != nil {
		t.Errorf("expected error: %s\n", err)
	}

	if cfg.Token != val["token"] {
		t.Error("token was not parsed")
	}
  if len(cfg.Admins) != 1 {
    t.Error("admins was not parsed")
  }
  if len(cfg.Mods) != 1 {
    t.Error("mods was not parsed")
  }
}

func TestParseConfig_returns_error_if_type_is_not_discord(t *testing.T) {
	val := make(map[string]interface{})
	val["type"] = "irc"

	_, err := ParseConfig(val)
	if err != ErrIncorrectType {
		t.Errorf("expected error: %s\n", err)
	}
}
