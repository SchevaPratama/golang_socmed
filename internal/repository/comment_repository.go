package repository

import (
	"golang_socmed/internal/entity"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CommentRepository struct {
	DB *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) List(postIds []string) ([]entity.Comment, error) {

	if len(postIds) < 1 {
		return nil, nil
	}
	query := "SELECT c.id, c.comment, c.created_at, c.user_id, c.post_id, u.name as user_name, u.friends as user_friends, u.createdat as user_created_at FROM comments as c LEFT JOIN users as u ON u.id::text = c.user_id::text WHERE c.post_id IN (" + strings.Join(postIds, ",") + ") ORDER BY c.created_at DESC"

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Slice to hold the fetched products
	var comments []entity.Comment

	// Loop through the rows and scan each product into the slice
	for rows.Next() {
		var comment entity.Comment
		// Use pq.Array to scan the Tags column into the product.Tags slice
		err := rows.Scan(&comment.ID, &comment.Comment, &comment.CreatedAt, &comment.UserId, &comment.PostId, &comment.UserName, pq.Array(&comment.UserFriends), &comment.UserCreatedAt)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) Create(request *entity.Comment) error {
	query := `INSERT INTO comments VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, request.ID, request.PostId, request.UserId, request.Comment)
	if err != nil {
		return err
	}

	return nil
}
