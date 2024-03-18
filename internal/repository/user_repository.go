package repository

import (
	"golang_socmed/internal/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(request *entity.User) error {
	query := `INSERT INTO users (name, username, password) VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(query, request.Name, request.Username, request.Password)
	return err
}

func (r *UserRepository) GetByUsername(request *entity.User) error {
	query := `SELECT * from users where username = $1`

	err := r.DB.Get(request, query, request.Username)
	return err
}

func (r *UserRepository) GetById(request *entity.User) error {
	query := `SELECT * from users where username = $1`

	err := r.DB.Get(request, query, request.ID)
	return err
}
