package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Domain string `json:"domain"`
	Dms    struct {
		DeviceStore string `json:"device_store"`
		DmsStore    string `json:"dms_store"`
		Endpoint    string `json:"endpoint"`
		CN          string `json:"common_name"`
		C           string `json:"country"`
		L           string `json:"locality"`
		O           string `json:"organization"`
		OU          string `json:"organization_unit"`
		ST          string `json:"state"`
		Name        string `json:"dms_name"`
	} `json:"dms"`
	DevManager struct {
		DevCrt  string `json:"devcrt"`
		DevAddr string `json:"addr"`
	} `json:"devmanager"`
	Auth struct {
		Endpoint string `json:"endpoint"`
	} `json:"auth"`
}

func NewConfig() (Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
