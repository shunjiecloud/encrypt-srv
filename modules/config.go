package modules

import (
	"os"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)

type KeyPairConfig struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func setupConfig() {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		panic("CONFIG_FILE_PATH is error")
	}
	if err := config.Load(file.NewSource(
		file.WithPath(configFilePath),
	)); err != nil {
		panic(err)
	}
}

func init() {
	setupConfig()
}
