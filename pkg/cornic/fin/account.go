package fin

import (
	"time"

	gdax "github.com/preichenberger/go-gdax"
)

// Account sructure has all the functions to manage an object
type Account struct {
	Name        string
	Currency    *Currency
	Balance     float64
	Available   float64
	LastUpdated time.Time
	ID          string
}

// NewAccount constructs an Account
func NewAccount(accountName string, accountCurrency *Currency, accountBalance float64, accountAvailable float64, accountID string) *Account {
	return &Account{Name: accountName, Currency: accountCurrency, Balance: accountBalance, Available: accountAvailable, LastUpdated: time.Now(), ID: accountID}
}

// updateAccountTime updates the lastUpdated time in the Account
func (a *Account) updateAccountTime() int {
	a.LastUpdated = time.Now()
	return 0
}

// SetBalance updates Account.Balance
func (a *Account) SetBalance(newBalance float64) int {
	a.Balance = newBalance
	a.updateAccountTime()
	return 0
}

// SetAvailable updates Account.Balance
func (a *Account) SetAvailable(newAvailable float64) int {
	a.Available = newAvailable
	a.updateAccountTime()
	return 0
}

// BuildAccountsList builds a list with all the used accounts
func BuildAccountsList(currencyList []*Currency, accounts []gdax.Account) map[string]*Account {
	allAccounts := make(map[string]*Account)
	for _, a := range accounts {
		for _, c := range currencyList {
			if c.Name == a.Currency {
				allAccounts[a.Currency] = NewAccount(a.Currency+" Account", c, a.Balance, a.Available, a.Id)
			}
		}
	}
	return allAccounts
}
