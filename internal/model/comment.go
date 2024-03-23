package model

type CommentResponse struct {
	Comment   string
	Creator   FriendResponse
	CreatedAt string
}

type CommentRequest struct {
	Comment string `json:"comment" validate:"required,min=2,max=500"`
	PostId  string `json:"postId" validate:"required"`
}
