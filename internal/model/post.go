package model

type PostFilter struct {
	Search     string   `query:"search" validate:"omitempty,alphanum"`
	SearchTags []string `query:"searchTag" validate:"omitempty,dive,required"`
	Limit      int      `query:"limit" validate:"number,min=1"`
	Offset     int      `query:"offset" validate:"number,min=0"`
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
	PostId   string            `json:"post_id"`
	Post     Post              `json:"post"`
	Creator  FriendResponse    `json:"creator"`
	Comments []CommentResponse `json:"comment"`
}

type PostRequest struct {
	PostInHtml string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags       []string `json:"tags" validate:"required,dive,required"`
}
