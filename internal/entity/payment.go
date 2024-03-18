package entity

type Payment struct {
	ID                   string
	BankAccountId        string `db:"bank_account_id"`
	ProductId            string `db:"product_id"`
	PaymentProofImageUrl string `db:"payment_proof_image_url"`
	Quantity             int
}

func (p *Payment) TableName() string {
	return "payments"
}
