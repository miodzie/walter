package rss

import (
	"github.com/maxatome/go-testdeep/helpers/tdsuite"
	"github.com/maxatome/go-testdeep/td"
	"testing"
)

// Fetch Feeds -> Notification -> Announcements -> Filters -> Deliver

// TODO
// - ThrottledMessenger Decorator
// - Update Subscription seen items on successful delivery?

type ProcessorSuite struct {
	processor  *Processor
	messenger  *StubMessenger
	repository *InMemRepository
	fetcher    *StubFetcher
}

func (p *ProcessorSuite) PreTest(t *td.T, testName string) error {
	p.repository = NewInMemRepo()
	p.fetcher = &StubFetcher{}
	p.messenger = &StubMessenger{}
	p.processor = NewProcessor(p.fetcher, p.repository, p.messenger)
	return nil
}

func (p *ProcessorSuite) TestDeliversNotifications(assert *td.T) {
	_ = p.processor.Process()

	assert.Len(p.messenger.delivered, 1)
}

func TestRunProcessorSuite(t *testing.T) {
	tdsuite.Run(t, new(ProcessorSuite))
}

///////////////////////////////////////////////////////

type StubMessenger struct {
	delivered []Announcement
}

func (m *StubMessenger) Deliver(announcements []Announcement) error {
	m.delivered = announcements
	return nil
}
