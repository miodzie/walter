package seras

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

var connectors map[string]BotParser

//go:embed config.toml
var DefaultConfig string

type Config struct {
	Mods    []string
	Bots    map[string]map[string]any
	Storage map[string]map[string]string
}

// BotParser intakes a map of config settings for a particular bot type.
type BotParser interface {
	Parse(map[string]any) (Bot, error)
}

func AddBotParser(name string, parser BotParser) error {
	if connectors == nil {
		connectors = make(map[string]BotParser)
	}
	if _, ok := connectors[name]; ok {
		return errors.New("connector already registered")
	}
	connectors[name] = parser
	return nil
}

func ParseToml(file string) (*Config, error) {
	var c Config
	f, err := os.ReadFile(file)
	if err != nil {
		return &c, err
	}
	err = toml.Unmarshal(f, &c)
	if err != nil {
		return &c, err
	}

	return &c, nil
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}
