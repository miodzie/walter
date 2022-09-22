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

	return ListSubscriptionsResponse{
		Message: "Success.",
		Subscriptions: []ListSubSubscription{
			{Feed: "news", Channel: "#general"},
		},
		Error: nil,
	}
}
