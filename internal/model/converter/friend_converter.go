package converter

import (
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"log"
	"time"
)

func FriendConverter(user *entity.User) *model.FriendResponse {
	t, err := time.Parse(time.RFC3339, user.CreatedAt)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
	}
	return &model.FriendResponse{
		UserId:      user.ID,
		Name:        user.Name,
		FriendCount: len(user.Friends),
		CreatedAt:   t.Format(time.RFC3339),
	}
}
