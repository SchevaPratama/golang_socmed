package converter

import (
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"log"
	"time"
)

func PostConverter(post *entity.Post) *model.PostResponse {
	postCreatedAt, err := time.Parse(time.RFC3339Nano, post.CreatedAt)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
	}

	userCreatedAt, err := time.Parse(time.RFC3339Nano, post.UserCreatedAt)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
	}

	return &model.PostResponse{
		PostId: post.ID,
		Post: model.Post{
			PostInHtml: post.PostInHtml,
			Tags:       post.Tags,
			CreatedAt:  postCreatedAt.Format(time.RFC3339),
		},
		Creator: model.PostCreator{
			UserId:      post.UserId,
			Name:        post.UserName,
			FriendCount: len(post.UserFriends),
			CreatedAt:   userCreatedAt.Format(time.RFC3339),
		},
	}
}
