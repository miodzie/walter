package usecases

import (
	"github.com/miodzie/seras/mods/rss"
)

type ListSubscriptions struct {
	repository rss.Repository
}

func NewListSubscriptionsUseCase(repository rss.Repository) *ListSubscriptions {
	return &ListSubscriptions{repository: repository}
}

type ListSubscriptionsRequest struct {
	User     string
	Optional struct{ Channel string }
}

type ListSubSubscription struct {
	Feed     string
	Channel  string
	Keywords []string
}

type ListSubscriptionsResponse struct {
	Message       string
	Subscriptions []ListSubSubscription
	Error         error
}

func (useCase ListSubscriptions) Handle(request ListSubscriptionsRequest) ListSubscriptionsResponse {
	var lsubs []ListSubSubscription
	subs, err := useCase.repository.Subs(rss.SubSearchOpt{
		User:    request.User,
		Channel: request.Optional.Channel,
	})
	if err != nil {
		return ListSubscriptionsResponse{
			Message: "Failed to retrieve subscriptions.",
			Error:   err,
		}
	}

	for _, sub := range subs {
		lsubs = append(lsubs, ListSubSubscription{
			Feed:     sub.Feed.Name,
			Channel:  sub.Channel,
			Keywords: sub.KeywordsSlice(),
		})
	}

	return ListSubscriptionsResponse{
		Message:       "Success.",
		Subscriptions: lsubs,
		Error:         nil,
	}
}
