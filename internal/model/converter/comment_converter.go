package converter

import (
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"log"
	"time"
)

func CommentConverter(comment *entity.Comment) *model.CommentResponse {
	commentCreatedAt, err := time.Parse(time.RFC3339Nano, comment.CreatedAt)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
	}

	userCreatedAt, err := time.Parse(time.RFC3339Nano, comment.UserCreatedAt)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
	}

	return &model.CommentResponse{
		Comment:   comment.Comment,
		CreatedAt: commentCreatedAt.Format(time.RFC3339),
		Creator: model.FriendResponse{
			UserId:      comment.UserId,
			Name:        comment.UserName,
			FriendCount: len(comment.UserFriends),
			CreatedAt:   userCreatedAt.Format(time.RFC3339),
		},
	}
}
