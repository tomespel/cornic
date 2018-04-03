package main

import (
	"math"
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

	fin "../pkg/cornic/fin"
	cornicio "../pkg/cornic/io"
)

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func main() {
	//fmt.Println("Hello")
	//fmt.Println(cornic.TestStrategy())

	cfg, err := cornicio.LoadConfiguration("../configs/practice.config.json")
	if err != nil {
		println(err.Error())
	}

	client := gdax.NewClient(cfg.Exchange.Secret, cfg.Exchange.Key, cfg.Exchange.Passphrase)
	client.HttpClient = &http.Client{Timeout: 15 * time.Second}

	// Building currency portfolio
	allCurrencies := fin.BuildCurrenciesList(cfg)

	accounts, err := client.GetAccounts()
	allAccounts := fin.BuildAccountsList(allCurrencies, accounts)
	if err != nil {
		println(err.Error())
	}

	allOrders := make(map[string]*fin.Order)
	activeOrdersList := make([]string, 0)

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := gdax.Message{
		Type: "subscribe",
		Channels: []gdax.MessageChannel{
			gdax.MessageChannel{Name: "ticker", ProductIds: cfg.Trading.TradedProducts}},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	i := 0
	message := gdax.Message{}

	secureThreshold := float64(0.0)
	action := true
	lastOrderID := ""
	triggerLevel := 1
	triggerAccountsSync := 20
	minimumInvestmentRequirement := float64(15)
	fee := float64(0.0)

	buyPrice := float64(0.0)
	sellValue := float64(0.0)

	for true {
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		for _, currency := range allCurrencies {
			if (message.ProductId[:3] == currency.Name) || (message.ProductId[4:] == currency.Name) {
				currency.Update(message.ProductId, "buy", message.BestAsk)
				currency.Update(message.ProductId, "sell", message.BestBid)
			}
		}

		// Show Update
		if message.ProductId[:3] == "ETH" {
			println(allCurrencies[1].Name, ":", allCurrencies[1].Rates[0].Name, "=", allCurrencies[1].Rates[0].ComputeMid(), allCurrencies[1].Rates[0].Fiat)
		}

		// Trigger sale

		if i%triggerLevel == 0 && i > triggerLevel && message.ProductId[:3] == "ETH" {

			if lastValue > allCurrencies[1].Rates[0].BidValue || lastValue < allCurrencies[1].Rates[0].AskPrice {
				if lastValue < allCurrencies[1].Rates[0].AskPrice && lastExecution < allCurrencies[1].Rates[0].AskPrice*(1-fee) {
					action = true
					println(">>", "Increase", action)
				}
				if lastValue > allCurrencies[1].Rates[0].BidValue && lastExecution > allCurrencies[1].Rates[0].BidValue*(1+fee) {
					action = true
					println(">>", "Decrease", action)
				}
			} else {
				action = false
			}

			if lastValue/allCurrencies[1].Rates[0].ComputeMid() < 1-secureThreshold || lastValue/allCurrencies[1].Rates[0].ComputeMid() > 1+secureThreshold {
			} else {

				// EURO ACCOUNT

				buyPrice = Round(allCurrencies[1].Rates[0].BidValue-0.01, 0.01)

				for _, order := range allOrders {
					if order.Price != buyPrice && order.Side == "buy" {
						err := client.CancelOrder(order.ID)
						if err != nil {
							println("Updated order:", err.Error())
						} else {
							println("Updated order:", "Cancelled", order.Side, "order", order.ID)
						}
					}
				}

				if allAccounts["EUR"].Available > minimumInvestmentRequirement && action == true {

					println("- -", "buyPrice", buyPrice)

					order := gdax.Order{
						Price:     buyPrice,
						Size:      0.01,
						Side:      "buy",
						ProductId: "ETH-EUR",
						PostOnly:  true,
					}
					savedOrder, err := client.CreateOrder(&order)
					if err != nil {
						println(err.Error())
						println("Attempted to buy at", buyPrice)
					} else {
						println("Order:", "buy at", lastExecution, "- PostOnly:", savedOrder.PostOnly, "- ID:", savedOrder.Id)
						allAccounts["EUR"].SetAvailable(allAccounts["EUR"].Available - lastExecution*0.01)
						println("Available balance:", allAccounts["EUR"].Available, allAccounts["EUR"].Currency.Name)
						activeOrdersList = append(activeOrdersList, lastOrderID)
					}

				}

				// ETHEREUM

				sellValue = Round(allCurrencies[1].Rates[0].AskPrice+0.01, 0.01)

				for _, order := range allOrders {
					if order.Price != buyPrice && order.Side == "sell" {
						err := client.CancelOrder(order.ID)
						if err != nil {
							println("Updated order:", err.Error())
						} else {
							println("Updated order:", "Cancelled", order.Side, "order", order.ID)
						}
					}
				}

				if allAccounts["ETH"].Available > 0 && action == true {

					println("- -", "sellValue", sellValue)

					lastExecutionType = "sell"
					order := gdax.Order{
						Price:     sellValue,
						Size:      0.01,
						Side:      "sell",
						ProductId: "ETH-EUR",
						PostOnly:  true,
					}
					savedOrder, err := client.CreateOrder(&order)
					if err != nil {
						println(err.Error())
						println("Attempted to sell at", sellValue)
					} else {
						println("Order:", "sell at", lastExecution, "- PostOnly:", savedOrder.PostOnly, "- ID:", savedOrder.Id)
						activeOrdersList = append(activeOrdersList, lastOrderID)
						allAccounts["ETH"].SetAvailable(allAccounts["ETH"].Available - 0.01)
						println("Available balance:", allAccounts["ETH"].Available, allAccounts["ETH"].Currency.Name)
					}

				}

			}
			//println(" ")
			lastValue = allCurrencies[1].Rates[0].ComputeMid()
		}

		if i%triggerAccountsSync == 0 || i == 0 {

			println("Synchronizing accounts.")

			accounts, err := client.GetAccounts()
			allAccounts = fin.BuildAccountsList(allCurrencies, accounts)
			if err != nil {
				println(err.Error())
			}

			println(" ")
			for _, account := range allAccounts {
				println(account.Name, "(", account.ID, "):", account.Available, account.Currency.Name)
			}
			println(" ")

			println("Synchronizing orders.")

			orders := make([]gdax.Order, 0)
			for _, orderID := range activeOrdersList {
				newOrder, err := client.GetOrder(orderID)
				if err == nil {
					if newOrder.Status == "open" {
						orders = append(orders, newOrder)
					}
				}
			}

			allOrders = fin.BuildOrdersList(orders)
			if err != nil {
				println(err.Error())
			}

			println(" ")
			for _, order := range allOrders {
				println(order.Side, order.Size, order.ProductID[:3], "at", order.Price, order.ProductID[4:], "(", order.ID, "- postOnly:", order.PostOnly, ")")
			}
			println(" ")

		}

		if message.Type == "snapshot" {
			println("Snapshot")
			println(" ")
		}

		i++

	}

}
