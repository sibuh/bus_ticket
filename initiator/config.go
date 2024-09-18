package initiator

import (
	"log"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func InitConfig(path string, logger *slog.Logger) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read config", err)
		log.Fatal()
	}

}
