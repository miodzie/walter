// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"github.com/miodzie/walter/mods/rss/internal/internal/domain"
)

// Parser downloads a Feed.Url and translates it to a ParsedFeed to
// be checked by a Subscription.
type Parser interface {
	ParseURL(string) (*domain.ParsedFeed, error)
}

type StubParser struct {
	Parsed *domain.ParsedFeed
}

func (p *StubParser) ParseURL(url string) (*domain.ParsedFeed, error) {
	return p.Parsed, nil
}
