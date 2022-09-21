package usecases

import (
	"fmt"
	"github.com/miodzie/seras/mods/rss"
)

type Unsubscribe struct {
	repository rss.Repository
}

func NewUnsubscribeUseCase(repo rss.Repository) *Unsubscribe {
	return &Unsubscribe{repository: repo}
}

type UnsubscribeRequest struct {
	User     string
	Channel  string
	FeedName string
}

type UnsubscribeResponse struct {
	Message string
	Error   error
}

func (useCase Unsubscribe) Handle(request UnsubscribeRequest) UnsubscribeResponse {
	sub, err := useCase.repository.SubByUserFeedNameChannel(
		request.User, request.FeedName, request.Channel,
	)
	if err != nil || sub == nil {
		return UnsubscribeResponse{
			Message: "Failed to locate user subscription.",
			Error:   err,
		}
	}

	err = useCase.repository.RemoveSub(sub)
	if err != nil {
		return UnsubscribeResponse{
			Message: "Failed to unsubscribe.",
			Error:   err,
		}
	}

	return UnsubscribeResponse{
		Message: fmt.Sprintf("Successfully unsubscribed from `%s` feed.", request.FeedName),
		Error:   nil,
	}
}
