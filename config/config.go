package config

var Config struct {
	Origin        string `yaml:"origin"`
	ServerAddress string `yaml:"server_address"`
	Static        string `yaml:"static"`
}
