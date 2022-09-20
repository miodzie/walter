package usecases

import (
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

	return UnsubscribeResponse{
		Message: "Successfully unsubscribed from news feed.",
		Error:   nil,
	}
}
