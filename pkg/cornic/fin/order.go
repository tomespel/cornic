package fin

import gdax "github.com/preichenberger/go-gdax"

// Order encapsulates market data about an order
type Order struct {
	ID        string
	Size      float64
	Side      string
	ProductID string
	Price     float64
	PostOnly  bool
}

// NewOrder constructs an Order
func NewOrder(orderID string, orderSize float64, orderSide string, orderProductID string, orderPrice float64, orderPostOnly bool) *Order {
	return &Order{ID: orderID, Size: orderSize, Side: orderSide, ProductID: orderProductID, Price: orderPrice, PostOnly: orderPostOnly}
}

// BuildOrdersList builds a list with all the current orders
func BuildOrdersList(orders []gdax.Order) map[string]*Order {
	allOrders := make(map[string]*Order)
	for _, o := range orders {
		allOrders[o.Id] = NewOrder(o.Id, o.Size, o.Side, o.ProductId, o.Price, o.PostOnly)
	}
	return allOrders
}
