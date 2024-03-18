package entity

type BankAccount struct {
	ID       string
	Name     string
	Number   string
	UserId   string `db:"user_id"`
	BankName string `db:"bank_name"`
}

func (ba *BankAccount) TableName() string {
	return "bank_accounts"
}
