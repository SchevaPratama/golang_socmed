package entity

type Post struct {
	ID            string
	PostInHtml    string `db:"post_in_html"`
	Tags          []string
	CreatedAt     string   `db:"created_at"`
	UserId        string   `db:"user_id"`
	UserName      string   `db:"user_name"`
	UserFriends   []string `db:"user_friends"`
	UserCreatedAt string   `db:"user_created_at"`
}

func (prod *Post) TableName() string {
	return "posts"
}
