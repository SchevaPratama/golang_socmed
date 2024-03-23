package model

type CommentResponse struct {
	Comment   string         `json:"comment"`
	Creator   FriendResponse `json:"creator"`
	CreatedAt string         `json:"createdAt"`
}

type CommentRequest struct {
	Comment string `json:"comment" validate:"required,min=2,max=500"`
	PostId  string `json:"postId" validate:"required"`
}
