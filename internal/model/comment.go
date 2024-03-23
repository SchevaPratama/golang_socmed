package model

type CommentResponse struct {
	Comment   string
	Creator   FriendResponse
	CreatedAt string
}
