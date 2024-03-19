package repository

import (
	"golang_socmed/internal/entity"
	"strconv"

	"github.com/jmoiron/sqlx"
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
	query += `SELECT ` + credentialType + `, name, password from users where ` + credentialType + ` = $` + strconv.Itoa(len(filterValues)+1)
	filterValues = append(filterValues, credentialValue)

	err := r.DB.Get(request, query, filterValues...)
	return err
}

func (r *UserRepository) GetById(request *entity.User) error {
	query := `SELECT * from users where username = $1`

	err := r.DB.Get(request, query, request.ID)
	return err
}
