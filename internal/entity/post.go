package entity

type Post struct {
	ID           string
	PostInHtml   string `db:"post_in_html"`
	Tags         []string
	UserId       string `db:"user_id"`
	UserName     string `db:"user_name"`
	UserImageUrl string `db:"user_image_url"`
	CreatedAt    string `db:"created_at"`
}

func (prod *Post) TableName() string {
	return "posts"
}
