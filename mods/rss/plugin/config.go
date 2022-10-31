// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package plugin

import (
	"errors"
	"fmt"
	"github.com/miodzie/seras"
	"github.com/miodzie/seras/mods/rss"
	"github.com/miodzie/seras/mods/rss/parsers/decorators"
	"github.com/miodzie/seras/mods/rss/parsers/gofeed"
	"github.com/miodzie/seras/storage"
	"github.com/miodzie/seras/storage/sqlite"
	"strings"
)

var parsers map[string]rss.Parser
var formatters map[string]rss.Formatter
var storages map[string]func(database string) (rss.Repository, error)

func init() {
	parsers = make(map[string]rss.Parser)
	parsers["gofeed"] = gofeed.New()

	formatters = make(map[string]rss.Formatter)
	formatters["default"] = rss.DefaultFormatter{}
	formatters["minimal"] = rss.MinimalFormatter{}

	storages = make(map[string]func(database string) (rss.Repository, error))
	storages["memory"] = func(database string) (rss.Repository, error) {
		return rss.NewInMemRepo(), nil
	}
	storages["sqlite"] = func(database string) (rss.Repository, error) {
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
	ctx := Context{}
	var ok bool
	ctx.Parser, ok = parsers[c.Parser]
	if !ok {
		return nil, fmt.Errorf("unknown parser: `%s`", c.Parser)
	}
	if c.Striphtml {
		ctx.Parser = decorators.StripHtml(ctx.Parser)
	}
	ctx.Formatter, ok = formatters[c.Formatter]
	if !ok {
		return nil, fmt.Errorf("unknown formatter: `%s`", c.Formatter)
	}
	createRepository, ok := storages[c.Storage]
	if !ok {
		return nil, fmt.Errorf("unknown storage: `%s`", c.Storage)
	}
	var err error
	ctx.Repository, err = createRepository(c.Database)
	if err != nil {
		return nil, err
	}

	return New(ctx), nil
}

func (s *Config) FillStruct(m map[string]any) error {
	for k, v := range m {
		k = strings.Title(k)
		err := seras.SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

type ModFactory struct {
	Context Context
}

func (m ModFactory) Create(c any) (seras.Module, error) {
	var conf Config
	config, ok := c.(map[string]any)
	if !ok {
		return nil, errors.New("failed to type assert rss module config")
	}
	err := conf.FillStruct(config)
	if err != nil {
		return nil, err
	}

	return conf.CreateMod()
}
