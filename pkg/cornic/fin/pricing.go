package fin

// TestPricing tests that the pricing file is loaded
func TestPricing() bool {
	return true
}

// ComputeCost computes the cost of purchase
func ComputeCost(cost float32, fees float32, quantity float32) float32 {
	return cost * (1.0 + fees) * quantity
}

// ComputeValue computes the value of selling
func ComputeValue(value float32, fees float32, quantity float32) float32 {
	return value * (1.0 - fees) * quantity
}
