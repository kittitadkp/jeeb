package domain

import "time"

type TransactionType string

const (
	TransactionIncome  TransactionType = "income"
	TransactionExpense TransactionType = "expense"
)

type Transaction struct {
	ID        string          `bson:"_id,omitempty" json:"id"`
	UserID    string          `bson:"user_id" json:"user_id"`
	Type      TransactionType `bson:"type" json:"type"`
	Amount    float64         `bson:"amount" json:"amount"`
	Category  string          `bson:"category" json:"category"`
	Date      time.Time       `bson:"date" json:"date"`
	Notes     string          `bson:"notes" json:"notes"`
	CreatedAt time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time       `bson:"updated_at" json:"updated_at"`
}
