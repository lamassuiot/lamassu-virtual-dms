package configs

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port     string `required:"true" split_words:"true"`
	Protocol string `required:"true" split_words:"true"`

	HomePath string

	DeviceManagerCertFile string `split_words:"true"`
	DeviceManagerAddress  string `split_words:"true"`

	MutualTLSEnabled  bool   `split_words:"true"`
	MutualTLSClientCA string `split_words:"true"`
	BootstrapCert     string `split_words:"true"`
	CertFile          string `split_words:"true"`
	KeyFile           string `split_words:"true"`
}

func NewConfig(prefix string) (error, Config) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return err, Config{}
	}
	return nil, cfg
}
