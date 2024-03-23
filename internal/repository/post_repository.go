package repository

import (
	"fmt"
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostRepository struct {
	DB *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) List(filter *model.PostFilter, userId string) ([]entity.Post, error) {

	friendsQuery := "select id from users u where array['" + userId + "']::text[] <@ friends::text[]"
	// Execute the query
	friendsRows, err := r.DB.Query(friendsQuery)
	if err != nil {
		// Handle error
		return nil, err
	}

	var userIds = []string{"'" + userId + "'"}
	// Loop through the rows and scan each product into the slice
	for friendsRows.Next() {
		var friend entity.User
		// Use pq.Array to scan the Tags column into the user.Tags slice
		err := friendsRows.Scan(&friend.ID)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, "'"+friend.ID+"'")
	}

	query := "SELECT p.id, p.post_in_html, p.tags, p.created_at, p.user_id, u.name as user_name, u.friends as user_friends, u.createdat as user_created_at FROM posts as p LEFT JOIN users as u ON u.id::text = p.user_id::text WHERE p.user_id IN (" + strings.Join(userIds, ",") + ")"
	var filterValues []interface{}

	// Conditionally append filters
	if filter.Search != "" {
		query += `AND post_in_html ILIKE CONCAT('%', $1::TEXT, '%')`
		filterValues = append(filterValues, filter.Search)
	}

	if filter.SearchTags != nil && len(filter.SearchTags) > 0 {
		tags := strings.Join(filter.SearchTags, "','")
		query += " AND p.tags::text[] && ARRAY['" + tags + "']::text[]"
	}

	query += " ORDER BY p.created_at DESC"

	if filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Limit)
	}

	if filter.Offset != 0 {
		query += fmt.Sprintf(" OFFSET $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Offset)
	}

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
		err := rows.Scan(&post.ID, &post.PostInHtml, pq.Array(&post.Tags), &post.CreatedAt, &post.UserId, &post.UserName, pq.Array(&post.UserFriends), &post.UserCreatedAt)
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
	query := `INSERT INTO posts VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, request.ID, request.PostInHtml, pq.Array(request.Tags), request.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) Get(postId string, userId string, request *entity.Post) error {
	query := "SELECT id, user_id FROM posts WHERE id = $1"
	err := r.DB.Get(request, query, postId)
	if err != nil {
		return &fiber.Error{
			Code:    404,
			Message: "Not found",
		}
	}

	var isFriend bool
	err = r.DB.QueryRow(`SELECT $1 = ANY(friends) FROM users WHERE id = $2`, request.UserId, userId).Scan(&isFriend)
	if err != nil {
		return err
	}

	if !isFriend {
		return &fiber.Error{
			Code:    400,
			Message: "postId is not the user’s friend",
		}
	}

	return nil
}
