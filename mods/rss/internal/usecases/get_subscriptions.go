// Copyright 2022-present miodzie. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package usecases

type GetSubscriptions struct {
	repository Repository
}

func NewGetSubscriptions(repository Repository) *GetSubscriptions {
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
	subs, err := useCase.repository.Subs(SearchParams{
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
			Keywords: sub.KeyWords(),
		})
	}

	return GetSubscriptionResponse{Message: "Success.", Subscriptions: lsubs}, nil
}
