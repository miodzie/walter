// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package rss

import (
	"github.com/miodzie/walter/log"
	"sync"
)

// FeedProcessor iterates over each UserFeed, retrieves a Feed through the Fetcher,
// finds new UserFeed Items that haven't been seen by each User's Subscription to said UserFeed,
// then turns them into grouped, sendable Notifications.
type FeedProcessor struct {
	// Max notifications sent per channel per Process() call.
	ChannelLimit int
	repository   Repository
	parser       Fetcher
	sync.Mutex
}

func NewFeedProcessor(repo Repository, parser Fetcher) *FeedProcessor {
	return &FeedProcessor{
		repository:   repo,
		parser:       parser,
		ChannelLimit: 3,
	}
}

func (p *FeedProcessor) Process() (notes []*Notification, err error) {
	p.Lock()
	defer p.Unlock()
	feeds, err := p.repository.Feeds()
	if err != nil {
		return notes, err
	}

	// FetchFeeds -> Notification -> Announcements

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

func (p *FeedProcessor) getSubsAndParsedFeed(feed *UserFeed) (
	[]*Subscription, *Feed, error) {
	parsedFeed, err := p.parser.Fetch(feed.Url)
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
func (p *FeedProcessor) updateSubscriptions(subs []*Subscription) {
	for _, s := range subs {
		if err := p.repository.UpdateSub(s); err != nil {
			log.Error(err)
		}
	}
}
