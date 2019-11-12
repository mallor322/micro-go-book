package config

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/spf13/viper"

	"log"
	"os"
)

const (
	kConfigType = "CONFIG_TYPE"
)

var Logger *log.Logger
var KitLogger kitlog.Logger


func init() {
	Logger = log.New(os.Stderr, "", log.LstdFlags)
	viper.AutomaticEnv()

	KitLogger = kitlog.NewLogfmtLogger(os.Stderr)
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.DefaultTimestampUTC)
	KitLogger = kitlog.With(KitLogger, "caller", kitlog.DefaultCaller)


	initDefault()

}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}

