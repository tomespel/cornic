package main

import (
	"net/http"
	"strconv"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

	fin "../pkg/cornic/fin"
	cornicio "../pkg/cornic/io"
)

func main() {
	//fmt.Println("Hello")
	//fmt.Println(cornic.TestStrategy())

	cfg, err := cornicio.LoadConfiguration("../configs/practice.config.json")
	if err != nil {
		println(err.Error())
	}

	client := gdax.NewClient(cfg.Exchange.Secret, cfg.Exchange.Key, cfg.Exchange.Passphrase)
	client.HttpClient = &http.Client{Timeout: 15 * time.Second}

	accounts, err := client.GetAccounts()
	if err != nil {
		println(err.Error())
	}

	for _, a := range accounts {
		println(a.Currency, "(", a.Id, "):", a.Balance)
	}

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := gdax.Message{
		Type: "subscribe",
		Channels: []gdax.MessageChannel{
			gdax.MessageChannel{Name: "level2", ProductIds: cfg.Trading.TradedAssets}},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	// Building currency portfolio
	allCurrencies := make([]*fin.Currency, 0)
	for _, element := range cfg.Trading.TradedAssets {
		newCurrency := fin.NewCurrency(element[:3], element[4:])
		allCurrencies = append(allCurrencies, newCurrency)
	}

	i := 0
	message := gdax.Message{}

	for true {
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		i++
		if i%200 == 0 {
			println("=======")
			for _, element := range allCurrencies {
				println(element.Name, "=", element.ComputeMid(), element.Fiat)
			}
		}

		if len(message.Changes) > 0 {
			if len(message.Changes[0]) > 0 {
				for _, c := range allCurrencies {
					if (c.Name == message.ProductId[:3]) && (c.Fiat == message.ProductId[4:]) {
						value, _ := strconv.ParseFloat(message.Changes[0][1], 64)
						c.Update(message.Changes[0][0], value)
					}
				}
			}
		}

		if message.Type == "snapshot" {
			println("Got a snapshot")
		}
	}

}
