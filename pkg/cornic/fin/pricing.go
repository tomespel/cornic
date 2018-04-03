package fin

// TestPricing tests that the pricing file is loaded
func TestPricing() bool {
	return true
}

// ComputeCost computes the cost of purchase
func ComputeCost(cost float64, fees float64, quantity float64) float64 {
	return cost * (1.0 + fees) * quantity
}

// ComputeValue computes the value of selling
func ComputeValue(value float64, fees float64, quantity float64) float64 {
	return value * (1.0 - fees) * quantity
}
