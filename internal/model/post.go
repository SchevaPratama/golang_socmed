package model

type PostFilter struct {
	Search     string   `query:"search" validate:"omitempty"`
	SearchTags []string `query:"searchTag" validate:"omitempty,dive,required"`
	Limit      int      `query:"limit"`
	Offset     int      `query:"offset"`
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
	PostId   string
	Post     Post
	Creator  FriendResponse
	Comments []CommentResponse
}

type PostRequest struct {
	PostInHtml string   `json:"postInHtml"`
	Tags       []string `query:"tags" validate:"required,dive,required"`
}
