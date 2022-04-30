package memory

import "github.com/miodzie/seras/mods/rss"


type InMemListenerRepository struct {
	items []*rss.Listener
}

func (r *InMemListenerRepository) All() []*rss.Listener {
	return r.items
}

func (r *InMemListenerRepository) Save(listener *rss.Listener) error {
	r.items = append(r.items, listener)
	return nil
}
