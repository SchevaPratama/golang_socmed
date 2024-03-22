package converter

import (
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
)

func PostConverter(post *entity.Post) *model.PostResponse {
	return &model.PostResponse{
		PostId: post.ID,
		Post: &model.Post{
			PostInHtml: post.PostInHtml,
			Tags:       post.Tags,
			CreatedAt:  post.CreatedAt,
		},
	}
}
