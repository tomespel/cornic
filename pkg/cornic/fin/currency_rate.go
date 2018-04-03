package fin

import (
	"time"

	io "../io"
)

// Loading congifuration
var cfg, _ = io.LoadConfiguration("../configs/practice.config.json")

var requiredNumber = cfg.Calibration.RequiredValues

// CurrencyRate encapsulates market data about a currency rate
type CurrencyRate struct {
	Name         string
	Fiat         string
	BidValue     float64
	AskPrice     float64
	LastUpdated  time.Time
	LastBids     []float64
	LastAsks     []float64
	RecordedBids int
	RecordedAsks int
}

// NewCurrencyRate constructs a CurrencyRate
func NewCurrencyRate(name string, fiat string) *CurrencyRate {
	return &CurrencyRate{Name: name, Fiat: fiat, LastBids: make([]float64, requiredNumber), LastAsks: make([]float64, requiredNumber), LastUpdated: time.Now()}
}

// updateCurrencyRateTime updates the lastUpdated time in the CurrencyRate
func (c *CurrencyRate) updateCurrencyRateTime() int {
	c.LastUpdated = time.Now()
	return 0
}

// SetBid updates CurrencyRate.BidValue and CurrencyRate.LastBids
func (c *CurrencyRate) SetBid(newBidValue float64) int {
	c.BidValue = newBidValue
	//println("LastBidssize:", len(c.LastBids))
	c.LastBids = c.LastBids[len(c.LastBids)-requiredNumber+1:]
	c.LastBids = append(c.LastBids, newBidValue)
	c.updateCurrencyRateTime()
	c.RecordedBids++
	return 0
}

// SetAsk updates CurrencyRate.AskPrice and CurrencyRate.LastAsks
func (c *CurrencyRate) SetAsk(newAskValue float64) int {
	c.AskPrice = newAskValue
	c.LastAsks = c.LastAsks[len(c.LastAsks)-requiredNumber+1:]
	c.LastAsks = append(c.LastAsks, newAskValue)
	c.updateCurrencyRateTime()
	c.RecordedAsks++
	return 0
}

// Update updates currency from the l2update stream
func (c *CurrencyRate) Update(action string, price float64) int {
	if action == "buy" {
		c.SetAsk(price)
		return 0
	}
	if action == "sell" {
		c.SetBid(price)
		return 0
	}
	return 1
}

// ComputeBid computes bid value as an average
func (c CurrencyRate) ComputeBid() float64 {
	if c.RecordedBids > requiredNumber {
		var sum float64
		sum = 0
		for _, element := range c.LastBids {
			sum += element
		}
		return sum / float64(len(c.LastBids))
	}
	return c.BidValue
}

// ComputeAsk computes ask price as an average
func (c CurrencyRate) ComputeAsk() float64 {
	if c.RecordedAsks > requiredNumber {
		var sum float64
		sum = 0
		for _, element := range c.LastAsks {
			sum += element
		}
		return sum / float64(len(c.LastAsks))
	}
	return c.AskPrice
}

// ComputeMid computes mid price as an average
func (c CurrencyRate) ComputeMid() float64 {
	return (c.ComputeBid() + c.ComputeAsk()) / 2
}
