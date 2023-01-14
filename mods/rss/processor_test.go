package rss

import (
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

// Fetch Feeds -> Notification -> Announcements -> Filters -> Announce

// TODO
// - ThrottledMessenger Decorator
// - Update Subscription seen items on successful delivery?

type ProcessorSuite struct {
	processor  *Processor
	announcer  *StubAnnouncer
	repository *InMemRepository
	fetcher    *StubFetcher
}

func (p *ProcessorSuite) PreTest(t *td.T, testName string) error {
	p.repository = NewInMemRepo()
	p.fetcher = &StubFetcher{}
	p.announcer = &StubAnnouncer{}
	p.processor = NewProcessor(p.fetcher, p.repository, p.announcer)
	return nil
}

func (p *ProcessorSuite) TestDeliversNotifications(assert *td.T) {
	_ = p.processor.Process()

	assert.Len(p.announcer.delivered, 1)
}

func TestRunProcessorSuite(t *testing.T) {
	tdsuite.Run(t, new(ProcessorSuite))
}

///////////////////////////////////////////////////////

type StubAnnouncer struct {
	delivered []Announcement
}

func (m *StubAnnouncer) Announce(announcements []Announcement) error {
	m.delivered = announcements
	return nil
}
