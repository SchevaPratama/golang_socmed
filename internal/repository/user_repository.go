package repository

import (
	"errors"
	"fmt"
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(credentialType string, credentialValue string, request *entity.User) error {
	query := `INSERT INTO users`
	var queryParams []interface{}
	if credentialType == "email" {
		query += ` (email, name, password) VALUES ($1, $2, $3)`
		queryParams = append(queryParams, request.Email, request.Name, request.Password)
	}

	if credentialType == "phone" {
		query += ` (phone, name, password) VALUES ($1, $2, $3)`
		queryParams = append(queryParams, request.Phone, request.Name, request.Password)
	}

	_, err := r.DB.Exec(query, queryParams...)
	return err
}

func (r *UserRepository) GetUsers(filter *model.FriendFilter, userId string) ([]entity.User, error) {
	tx, _ := r.DB.Beginx()
	defer tx.Rollback()

	query := `SELECT id, name, friends, createdAt FROM users`
	var filterValues []interface{}

	// Conditionally append filters
	if *filter.Search != "" {
		query += ` WHERE name LIKE CONCAT('%', $1::TEXT, '%')`
		filterValues = append(filterValues, *filter.Search)
	}

	if *filter.OnlyFriend {
		if len(filterValues) > 0 {
			query += ` AND `
		} else {
			query += ` WHERE `
		}
		query += "ARRAY['" + userId + "']::text[] <@ friends::text[]"
	}

	// Add sorting if SortField and SortOrder are provided
	if *filter.SortBy != "" && *filter.OrderBy != "" {
		query += ` ORDER BY ` + *filter.SortBy + ` ` + *filter.OrderBy
	} else {
		query += ` ORDER BY createdAt DESC`
	}

	if *filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Limit)
	} else {
		query += " LIMIT 5"
	}

	if *filter.Offset != 0 {
		query += fmt.Sprintf(" OFFSET $%s", strconv.Itoa(len(filterValues)+1))
		filterValues = append(filterValues, filter.Offset)
	} else {
		query += " OFFSET 0"
	}

	log.Println(query)

	// Execute the query
	rows, err := r.DB.Query(query, filterValues...)
	if err != nil {
		// Handle error
		return nil, err
	}

	// Slice to hold the fetched products
	var users []entity.User

	// Loop through the rows and scan each product into the slice
	for rows.Next() {
		var user entity.User
		// Use pq.Array to scan the Tags column into the user.Tags slice
		err := rows.Scan(&user.ID, &user.Name, pq.Array(&user.Friends), &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// if err := query.Find(&product).Error; err != nil {
	// 	return nil, err
	// }

	return users, nil
}

func (r *UserRepository) GetByEmailOrPhone(credentialType string, credentialValue string, request *entity.User) error {
	var query string
	var filterValues []interface{}
	// if credentialType == "email" {
	// 	query += `SELECT` + credentialType + `, name from users where ` + credentialType + ` = ` + credentialValue
	// 	filterValues = append(filterValues, credentialValue)
	// }

	// if credentialType == "phone" {
	// 	query += `SELECT phone, name from users where phone = $` + strconv.Itoa(len(filterValues)+1)
	// 	filterValues = append(filterValues, credentialValue)
	// }
	query += `SELECT id, ` + credentialType + `, name, password from users where ` + credentialType + ` = $` + strconv.Itoa(len(filterValues)+1)
	filterValues = append(filterValues, credentialValue)

	err := r.DB.Get(request, query, filterValues...)
	return err
}

func (r *UserRepository) AddFriend(friendId string, userId string) error {
	query := `UPDATE users SET friends = array_append(friends, $1) WHERE id = $2`
	// query2 := `UPDATE users SET friends = array_append(friends, $1) WHERE id = $2`
	var exists bool
	err3 := r.DB.QueryRow(`SELECT $1 = ANY(friends) FROM users WHERE id = $2`, friendId, userId).Scan(&exists)
	if err3 != nil {
		return err3
	}

	log.Println(exists)

	if exists {
		return errors.New("User Already Friend")
	} else {
		_, err := r.DB.Exec(query, friendId, userId)
		if err != nil {
			return err
		}

		_, err2 := r.DB.Exec(query, userId, friendId)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func (r *UserRepository) DeleteFriend(friendId string, userId string) error {
	query := `UPDATE users SET friends = array_remove(friends, $1) WHERE id = $2`
	// query2 := `UPDATE users SET friends = array_append(friends, $1) WHERE id = $2`
	var exists bool
	err3 := r.DB.QueryRow(`SELECT $1 = ANY(friends) FROM users WHERE id = $2`, friendId, userId).Scan(&exists)
	if err3 != nil {
		return err3
	}

	if !exists {
		return errors.New("Friend Not Found")
	} else {
		_, err := r.DB.Exec(query, friendId, userId)
		if err != nil {
			return err
		}

		_, err2 := r.DB.Exec(query, userId, friendId)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func (r *UserRepository) LinkPhoneEmail(types string, value string, userId string) error {
	var query string
	// var filterValues []interface{}
	// if credentialType == "email" {
	// 	query += `SELECT` + credentialType + `, name from users where ` + credentialType + ` = ` + credentialValue
	// 	filterValues = append(filterValues, credentialValue)
	// }

	// if credentialType == "phone" {
	// 	query += `SELECT phone, name from users where phone = $` + strconv.Itoa(len(filterValues)+1)
	// 	filterValues = append(filterValues, credentialValue)
	// }
	query += `UPDATE users SET ` + types + ` = $1 WHERE id = $2`
	log.Println(query)
	log.Println(types)
	log.Println(value)

	_, err := r.DB.Exec(query, value, userId)
	return err
}

func (r *UserRepository) GetById(request *model.FriendRequest) (entity.User, error) {
	query := `SELECT id, email, phone, name, friends, createdAt FROM users WHERE id = $1`

	userData := entity.User{}

	// Execute the query
	rows, err := r.DB.Query(query, request.UserId)
	if err != nil {
		// Handle error
		return userData, err
	}

	// Loop through the rows and scan each product into the slice
	for rows.Next() {
		var user entity.User
		// Use pq.Array to scan the Tags column into the product.Tags slice
		err := rows.Scan(&user.ID, &user.Email, &user.Phone, &user.Name, pq.Array(&user.Friends), &user.CreatedAt)
		if err != nil {
			return userData, err
		}
		userData = user
	}

	if !userData.Phone.Valid {
		userData.Phone.String = "" // Set to empty string or any default value
	}

	if !userData.Email.Valid {
		userData.Email.String = "" // Set to empty string or any default value
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return userData, err
	}
	// Check for no results
	if userData.Name == "" {
		return userData, errors.New("No Data Found")
	}

	return userData, nil
}
