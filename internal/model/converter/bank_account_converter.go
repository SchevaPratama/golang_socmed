package converter

import (
	"golang_socmed/internal/entity"
	"golang_socmed/internal/model"
)

func BankAccountConverter(bankAccount *entity.BankAccount) *model.BankAccountResponse {
	return &model.BankAccountResponse{
		BankAccountId:     bankAccount.ID,
		BankName:          bankAccount.BankName,
		BankAccountName:   bankAccount.Name,
		BankAccountNumber: bankAccount.Number,
	}
}
