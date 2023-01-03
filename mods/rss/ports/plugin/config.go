// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"errors"
	"fmt"
	"github.com/miodzie/walter"
	domain "github.com/miodzie/walter/mods/rss/internal/usecases"
	"github.com/miodzie/walter/mods/rss/internal/usecases/adapters/parsers/decorators"
	"github.com/miodzie/walter/mods/rss/internal/usecases/adapters/parsers/gofeed"
	"github.com/miodzie/walter/mods/rss/internal/usecases/adapters/storage/sqlite"
	"github.com/miodzie/walter/storage"
	"strings"
)

var parsers map[string]domain.Parser
var formatters map[string]domain.Formatter
var storages map[string]func(database string) (domain.Repository, error)

func init() {
	parsers = make(map[string]domain.Parser)
	parsers["gofeed"] = gofeed.New()

	formatters = make(map[string]domain.Formatter)
	formatters["default"] = domain.DefaultFormatter{}
	formatters["minimal"] = domain.MinimalFormatter{}

	storages = make(map[string]func(database string) (domain.Repository, error))
	storages["memory"] = func(database string) (domain.Repository, error) {
		return domain.NewInMemRepo(), nil
	}
	storages["sqlite"] = func(database string) (domain.Repository, error) {
		db, err := storage.Get(database)
		if err != nil {
			return nil, err
		}
		return sqlite.NewRssRepository(db), nil
	}
}

type Config struct {
	Parser    string
	Striphtml bool
	Formatter string
	Storage   string
	Database  string
}

func (c *Config) CreateMod() (*RssMod, error) {
	services := Services{}
	var ok bool
	services.Parser, ok = parsers[c.Parser]
	if !ok {
		return nil, fmt.Errorf("unknown parser: `%s`", c.Parser)
	}
	if c.Striphtml {
		services.Parser = decorators.StripHtml(services.Parser)
	}
	services.Formatter, ok = formatters[c.Formatter]
	if !ok {
		return nil, fmt.Errorf("unknown formatter: `%s`", c.Formatter)
	}
	createRepository, ok := storages[c.Storage]
	if !ok {
		return nil, fmt.Errorf("unknown storage: `%s`", c.Storage)
	}
	var err error
	services.Repository, err = createRepository(c.Database)
	if err != nil {
		return nil, err
	}

	return New(services), nil
}

func (s *Config) FillStruct(m map[string]any) error {
	for k, v := range m {
		k = strings.Title(k)
		err := walter.SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

type ModFactory struct {
	Context Services
}

func (m ModFactory) Create(c any) (walter.Module, error) {
	var conf Config
	config, ok := c.(map[string]any)
	if !ok {
		return nil, errors.New("failed to type assert entitiy module config")
	}
	err := conf.FillStruct(config)
	if err != nil {
		return nil, err
	}

	return conf.CreateMod()
}
