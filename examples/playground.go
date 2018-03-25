package main

import (
	"net/http"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

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

	//println(len(cfg.Exchange.TradedMarkets))
	//println(cfg.Exchange.TradedMarkets)

	subscribe := gdax.Message{
		Type: "subscribe",
		Channels: []gdax.MessageChannel{
			gdax.MessageChannel{
				Name: "level2",
				ProductIds: []string{
					//"ETH-BTC",
					"ETH-BTC",
				},
			},
		},
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}
	i := 0
	message := gdax.Message{}
	for true {
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
			break
		}

		i++
		println(message.Type, message.ProductId, message.Price, i)
		if len(message.Changes) > 0 {
			if len(message.Changes[0]) > 0 {
				println(message.Type, message.ProductId, message.Changes[0][1], i)
			}
		}

		if message.Type == "snapshot" {
			println("Got a snapshot")
		}
	}

}
