package config

type ClientConfig struct {
	ServerURL  string
	ServerPort uint
}

type ServerConfig struct {
	Port      uint
	PublicURL string
	Workdir   string
	Salt      string
}

type Config struct {
	Passphrase string
	Client     ClientConfig
	Server     ServerConfig
}
