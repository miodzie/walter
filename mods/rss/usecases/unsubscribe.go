package usecases

import (
	"fmt"
	"github.com/miodzie/seras/mods/rss"
)

type Unsubscribe struct {
	repository rss.Repository
}

func NewUnsubscribe(repo rss.Repository) *Unsubscribe {
	return &Unsubscribe{repository: repo}
}

type UnsubscribeRequest struct {
	User     string
	Channel  string
	FeedName string
}

type UnsubscribeResponse struct {
	Message string
}

func (useCase Unsubscribe) Unsubscribe(
	request UnsubscribeRequest) (UnsubscribeResponse, error) {
	subs, err := useCase.repository.Subs(rss.SubSearchOpt{
		User:     request.User,
		FeedName: request.FeedName,
		Channel:  request.Channel,
	})

	if err != nil || len(subs) != 1 {
		return UnsubscribeResponse{Message: "Failed to locate user subscription."}, err
	}
	err = useCase.repository.RemoveSub(subs[0])
	if err != nil {
		return UnsubscribeResponse{Message: "Failed to unsubscribe."}, err
	}

	return UnsubscribeResponse{
		Message: fmt.Sprintf("Successfully unsubscribed from `%s` feed.", request.FeedName),
	}, nil
}
