# Finance Module

## Domain

```go
type Transaction struct {
    ID        string
    UserID    string
    Amount    float64
    Type      TransactionType  // income, expense
    Category  string
    Date      time.Time
    Notes     string
}
```

## Use Cases

- AddTransaction
- ListTransactions
- GetSpendingByCategory
- GetMonthlyReport
- SetBudgetAlert

## Categories

Income: salary, freelance, investment
Expense: food, transport, utilities, entertainment
