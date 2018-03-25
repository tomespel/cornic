package io

import (
	"encoding/json"
	"os"
)

// Config encapsulates config option
type Config struct {
	Owner struct {
		Name string
	}
	Exchange struct {
		Key        string
		Secret     string
		Passphrase string
	}
	Trading struct {
		Fees            float32
		RiskAversion    float32
		BuySensitivity  float32
		SellSensitivity float32
		ActionTick      int
	}
	ROI struct {
		PaymentPeriod int
		ProfitRate    float32
	}
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
