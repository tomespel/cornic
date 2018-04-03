package fin

import(
  io "../io"
)

// Currency encapsulates market data about a currency
type Currency struct{
  Name string
  Rates []*CurrencyRate
}

// NewCurrency constructs a Currency
func NewCurrency(name string) *Currency{
	return &Currency{Name: name}
}

// GetRate returns the CurrencyRate object
func (c Currency) GetRate(name string, fiat string) *CurrencyRate{
  for _, element := range c.Rates{
    if (element.Name==name) && (element.Fiat==fiat){
      return element
    }
  }
  return nil
}

// AddRate adds rates in pair for a given currency
func (c *Currency) AddRate(rateID string) int{
  c.Rates = append(c.Rates, NewCurrencyRate(rateID[:3], rateID[4:]))
  c.Rates = append(c.Rates, NewCurrencyRate(rateID[4:], rateID[:3]))
  return 0
}

// Update updates the rates from the l2update stream
func (c *Currency) Update(rateID string, action string, price float64) int {
  if action == "buy" {
		c.GetRate(rateID[:3], rateID[4:]).SetAsk(price)
    c.GetRate(rateID[4:], rateID[:3]).SetBid(1/price)
		return 0
	}
	if action == "sell" {
    c.GetRate(rateID[:3], rateID[4:]).SetBid(price)
    c.GetRate(rateID[4:], rateID[:3]).SetAsk(1/price)
		return 0
	}
  return 1
}

// BuildCurrenciesList builds a list with all the currencies
func BuildCurrenciesList(cfg io.Config) []*Currency{
  allCurrencies := make([]*Currency, 0)
  for _, asset := range cfg.Trading.TradedAssets {
    newCurrency := NewCurrency(asset)
    for _, product := range cfg.Trading.TradedProducts {
      if (product[:3] == asset) || (product[4:] == asset) {
        newCurrency.AddRate(product)
      }
    }
    allCurrencies = append(allCurrencies, newCurrency)
  }
  return allCurrencies
}
