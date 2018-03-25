package fin

import "time"

// Account sructure has all the functions to manage an object
type Account struct {
	Name        string
	Currency    string
	Balance     float32
	Available   float32
	LastUpdated time.Time
	ID          string
}

// NewAccount constructs an Account
func NewAccount(accountName string, accountCurrency string, accountBalance float32, accountAvailable float32, accountID string) *Account {
	return &Account{Name: accountName, Currency: accountCurrency, Balance: accountBalance, Available: accountAvailable, LastUpdated: time.Now(), ID: accountID}
}

// updateAccountTime updates the lastUpdated time in the Account
func (a *Account) updateAccountTime() int {
	a.LastUpdated = time.Now()
	return 0
}

// SetBalance updates Account.Balance
func (a *Account) SetBalance(newBalance float32) int {
	a.Balance = newBalance
	a.updateAccountTime()
	return 0
}

// SetAvailable updates Account.Balance
func (a *Account) SetAvailable(newAvailable float32) int {
	a.Available = newAvailable
	a.updateAccountTime()
	return 0
}
