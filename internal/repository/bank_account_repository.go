package repository

import (
	"fmt"
	"golang_socmed/internal/entity"

	"github.com/jmoiron/sqlx"
)

type BankAccountRepository struct {
	DB *sqlx.DB
}

func NewBankAccountRepository(db *sqlx.DB) *BankAccountRepository {
	return &BankAccountRepository{DB: db}
}

func (r *BankAccountRepository) List(userId string) ([]entity.BankAccount, error) {
	bankAccounts := []entity.BankAccount{}

	query := `SELECT id, name, number, bank_name, user_id FROM bank_accounts WHERE user_id = $1`
	err := r.DB.Select(&bankAccounts, query, userId)
	if err != nil {
		return nil, err
	}

	return bankAccounts, nil
}

func (r *BankAccountRepository) Get(id string, bankAccount *entity.BankAccount, userId string) error {

	query := `SELECT id, name, number, bank_name, user_id FROM bank_accounts WHERE id = $1 AND user_id = $2`
	err := r.DB.Get(bankAccount, query, id, userId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *BankAccountRepository) Create(request *entity.BankAccount) error {
	query := `INSERT INTO bank_accounts VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, request.ID, request.Name, request.Number, request.BankName, request.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *BankAccountRepository) Update(id string, request *entity.BankAccount) error {
	query := `UPDATE bank_accounts SET name = $1, number = $2, bank_name = $3 WHERE id = $4`
	_, err := r.DB.Exec(query, request.Name, request.Number, request.BankName, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *BankAccountRepository) Delete(id string) error {
	query := `DELETE FROM bank_accounts WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
