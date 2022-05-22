package discord

import "testing"

func TestFromConfig_parses_into_Config_struct(t *testing.T) {
	val := make(map[string]interface{})
	val["type"] = "discord"
	val["token"] = "my_token"
	admins := make(map[string]bool)
	admins["alice"] = true
	val["admins"] = admins
	val["mods"] = []string{"myplugin"}

	cfg, err := FromConfig(val)
	if err != nil {
		t.Errorf("expected error: %s\n", err)
	}

	if cfg.Token != val["token"] {
		t.Error("token was not parsed")
	}
	_, ok := cfg.Admins["alice"]
	if !ok {
		t.Error("admins was not parsed")
	}
  if len(cfg.Mods) != 1 {
    t.Error("mods was not parsed")
  }
}

func TestFromConfig_returns_error_if_type_is_not_discord(t *testing.T) {
	val := make(map[string]interface{})
	val["type"] = "irc"

	_, err := FromConfig(val)
	if err != ErrIncorrectType {
		t.Errorf("expected error: %s\n", err)
	}
}
