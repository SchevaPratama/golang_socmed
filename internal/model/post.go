package model

type PostFilter struct {
	Search     string   `query:"search" validate:"omitempty"`
	SearchTags []string `query:"searchTag" validate:"omitempty,dive,required"`
	Limit      int      `query:"limit" validate:"min=1"`
	Offset     string   `query:"offset" validate:"number,gte=0"`
}

type Post struct {
	PostInHtml string   `json:"postInHtml"`
	Tags       []string `json:"tags"`
	CreatedAt  string   `json:"createdAt"`
}

// type PostCreator struct {
// 	UserId      string `json:"userId"`
// 	Name        string
// 	ImageUrl    string
// 	FriendCount int
// 	CreatedAt   string `json:"createdAt"`
// }

type PostResponse struct {
	PostId   string            `json:"postId"`
	Post     Post              `json:"post"`
	Creator  FriendResponse    `json:"creator"`
	Comments []CommentResponse `json:"comments"`
}

type PostRequest struct {
	PostInHtml string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags       []string `json:"tags" validate:"required,dive,required"`
}
