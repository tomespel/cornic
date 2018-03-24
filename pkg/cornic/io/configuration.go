package io

import (
	"encoding/json"
	"os"
)

// Config encapsulates config option
type Config struct {
	Owner struct {
		Name string `json: "name"`
	} `json: "owner"`
	Exchange struct {
		Key        string `json: "key"`
		Secret     string `json: "secret"`
		Passphrase string `json: "passphrase"`
	} `json: "exchange"`
	Trading struct {
		Fees            float32 `json: "fees"`
		RiskAversion    float32 `json: "roskAversion"`
		BuySensitivity  float32 `json: "buySensitivity"`
		SellSensitivity float32 `json: "sellSensivity"`
	} `json: "trading"`
	ROI struct {
		PaymentPeriod int     `json: "paymentPeriod"`
		ProfitRate    float32 `json: "profitRate"`
	} `json: "roi"`
}

// LoadConfiguration loads configFile
func LoadConfiguration(fileName string) (Config, error) {
	var config Config
	configFile, err := os.Open(fileName)
	defer configFile.Close()
	if err != nil {
		return config, nil
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}
