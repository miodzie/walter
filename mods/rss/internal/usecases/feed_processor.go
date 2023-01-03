// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

import (
	"github.com/miodzie/walter/log"
	"github.com/miodzie/walter/mods/rss/internal/internal/domain"
	"sync"
)

// FeedProcessor iterates over each Feed, retrieves a ParsedFeed through the Parser,
// finds new Feed Items that haven't been seen by each User's Subscription to said Feed,
// then turns them into grouped, sendable Notifications.
type FeedProcessor struct {
	// Max notifications sent per channel per Process() call.
	ChannelLimit int
	repository   Repository
	parser       Parser
	sync.Mutex
}

func NewProcessor(repo Repository, parser Parser) *FeedProcessor {
	return &FeedProcessor{
		repository:   repo,
		parser:       parser,
		ChannelLimit: 3,
	}
}

func (p *FeedProcessor) Process() (notes []*domain.Notification, err error) {
	p.Lock()
	defer p.Unlock()
	feeds, err := p.repository.Feeds()
	if err != nil {
		return notes, err
	}

	sorter := newNoteSorter(p.ChannelLimit)
	for _, feed := range feeds {
		subs, parsedFeed, err2 := p.getSubsAndParsedFeed(feed)
		if err2 != nil {
			return notes, err2
		}
		sorted := sorter.sort(subs, parsedFeed)
		notes = append(notes, sorted...)
		p.updateSubscriptions(subs)
	}

	return notes, nil
}

func (p *FeedProcessor) getSubsAndParsedFeed(feed *domain.Feed) (
	[]*domain.Subscription, *domain.ParsedFeed, error) {
	parsedFeed, err := p.parser.ParseURL(feed.Url)
	if err != nil {
		return nil, nil, err
	}
	subs, err := p.repository.Subs(SearchParams{FeedId: feed.Id})
	if err != nil {
		return nil, nil, err
	}

	return subs, parsedFeed, nil
}

// TODO: Write a test for this.
// Hard to write a test for this when repositories use pointers for everything.
func (p *FeedProcessor) updateSubscriptions(subs []*domain.Subscription) {
	for _, s := range subs {
		if err := p.repository.UpdateSub(s); err != nil {
			log.Error(err)
		}
	}
}
