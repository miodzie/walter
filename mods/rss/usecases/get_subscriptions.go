package usecases

import (
	"github.com/miodzie/seras/mods/rss"
)

type GetSubscriptions struct {
	repository rss.Repository
}

func NewGetSubscriptions(repository rss.Repository) *GetSubscriptions {
	return &GetSubscriptions{repository: repository}
}

type GetSubscriptionsRequest struct {
	User     string
	Optional struct{ Channel string }
}

type GetSubscriptionsSub struct {
	Feed     string
	Channel  string
	Keywords []string
}

type GetSubscriptionResponse struct {
	Message       string
	Subscriptions []GetSubscriptionsSub
}

func (useCase GetSubscriptions) Get(request GetSubscriptionsRequest) (GetSubscriptionResponse, error) {
	var lsubs []GetSubscriptionsSub
	subs, err := useCase.repository.Subs(rss.SubSearchOpt{
		User:    request.User,
		Channel: request.Optional.Channel,
	})
	if err != nil {
		return GetSubscriptionResponse{Message: "Failed to retrieve subscriptions."}, err
	}

	for _, sub := range subs {
		lsubs = append(lsubs, GetSubscriptionsSub{
			Feed:     sub.Feed.Name,
			Channel:  sub.Channel,
			Keywords: sub.KeywordsSlice(),
		})
	}

	return GetSubscriptionResponse{Message: "Success.", Subscriptions: lsubs}, nil
}