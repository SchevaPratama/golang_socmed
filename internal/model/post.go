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

type PostResponse struct {
	PostId string
	Post   *Post
}