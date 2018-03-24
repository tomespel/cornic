package main

import (
	"net/http"
	"time"

	gdax "github.com/preichenberger/go-gdax"

	cornicio "../pkg/cornic/io"
)

func main() {
	//fmt.Println("Hello")
	//fmt.Println(cornic.TestStrategy())

	cfg, err := cornicio.LoadConfiguration("../configs/practice.config.yml")
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

}
