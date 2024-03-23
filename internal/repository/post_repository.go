package repository

import (
	"fmt"
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostRepository struct {
	DB *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) List(filter *model.PostFilter) ([]entity.Post, error) {

	query := `SELECT p.id, p.post_in_html, p.tags, p.created_at FROM posts as p LEFT JOIN users as u ON u.id::text = p.user_id::text`
	var filterValues []interface{}

	// Conditionally append filters
	if filter.Search != "" {
		query += ` WHERE post_in_html LIKE CONCAT('%', $1::TEXT, '%')`
		filterValues = append(filterValues, filter.Search)
	}

	if filter.SearchTags != nil && len(filter.SearchTags) > 0 {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}

		tags := strings.Join(filter.SearchTags, "','")
		query += "p.tags::text[] && ARRAY['" + tags + "']::text[]"
	}

	if filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Limit)
	}

	if filter.Offset != 0 {
		query += fmt.Sprintf(" OFFSET $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Offset)
	}

	fmt.Println(query)

	// Execute the query
	rows, err := r.DB.Query(query, filterValues...)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Slice to hold the fetched products
	var posts []entity.Post

	// Loop through the rows and scan each product into the slice
	for rows.Next() {
		var post entity.Post
		// Use pq.Array to scan the Tags column into the product.Tags slice
		err := rows.Scan(&post.ID, &post.PostInHtml, pq.Array(&post.Tags), &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) Create(request *entity.Post) error {
	fmt.Println(request)
	query := `INSERT INTO posts VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, request.ID, request.PostInHtml, pq.Array(request.Tags), request.UserId)
	if err != nil {
		return err
	}

	return nil
}
