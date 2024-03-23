package entity

type Comment struct {
	ID            string
	Comment       string
	PostId        string   `db:"post_id"`
	CreatedAt     string   `db:"created_at"`
	UserId        string   `db:"user_id"`
	UserName      string   `db:"user_name"`
	UserFriends   []string `db:"user_friends"`
	UserCreatedAt string   `db:"user_created_at"`
}

func (prod *Comment) TableName() string {
	return "comments"
}
