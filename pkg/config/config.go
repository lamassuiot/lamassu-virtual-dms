package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Dms struct {
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
		Username string `json:"username"`
		Password string `json:"Password"`
	} `json:"auth"`
}
type Token struct {
	AccesToken         string `json:"access_token"`
	Expires_in         string `json:"expires_in"`
	Policy             string `json:"not-before-policy"`
	Refresh_expires_in string `json:"refresh_expires_in"`
	Refresh_token      string `json:"refresh_token"`
	Scope              string `json:"scope"`
	Session_state      string `json:"session_state"`
	Token_type         string `json:"token_type"`
}
type Device struct {
	Id                      string                        `json:"id"`
	Alias                   string                        `json:"alias"`
	Description             string                        `json:"description"`
	Tags                    []string                      `json:"tags"`
	IconName                string                        `json:"iconName"`
	IconColor               string                        `json:"iconColor"`
	Status                  string                        `json:"status,omitempty"`
	DmsId                   string                        `json:"dms_id"`
	KeyMetadata             PrivateKeyMetadataWithStregth `json:"key_metadata"`
	Subject                 Subject                       `json:"subject"`
	CreationTimestamp       string                        `json:"creation_timestamp,omitempty"`
	CurrentCertSerialNumber string                        `json:"current_cert_serial_number"`
}
type PrivateKeyMetadataWithStregth struct {
	KeyType     string `json:"type"`
	KeyBits     int    `json:"bits"`
	KeyStrength string `json:"strength"`
}
type Subject struct {
	CN string `json:"common_name"`
	O  string `json:"organization"`
	OU string `json:"organization_unit"`
	C  string `json:"country"`
	ST string `json:"state"`
	L  string `json:"locality"`
}
type DMS struct {
	Id                    string             `json:"id"`
	Name                  string             `json:"name"`
	SerialNumber          string             `json:"serial_number,omitempty"`
	KeyMetadata           PrivateKeyMetadata `json:"key_metadata"`
	Status                string             `json:"status"`
	CsrBase64             string             `json:"csr,omitempty"`
	CerificateBase64      string             `json:"crt,omitempty"`
	Subject               Subject            `json:"subject"`
	AuthorizedCAs         []string           `json:"authorized_cas,omitempty"`
	CreationTimestamp     string             `json:"creation_timestamp,omitempty"`
	ModificationTimestamp string             `json:"modification_timestamp,omitempty"`
}
type PrivateKeyMetadata struct {
	KeyType string `json:"type"`
	KeyBits int    `json:"bits"`
}
type DmsCreationResponse struct {
	Dms     DMS    `json:"dms,omitempty"`
	PrivKey string `json:"priv_key,omitempty"`
	Err     error  `json:"err,omitempty"`
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
