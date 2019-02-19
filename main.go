package main

import (
	"fmt"
	"log"

	"github.com/ltran/priceit/rideshare"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")

	// Find and read the config file and handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	lyftcfg := viper.GetStringMapString("app.lyft")
	auth := rideshare.LyftAuth(nil, lyftcfg["username"], lyftcfg["password"])
	lyftEsts := rideshare.LyftCostEstimate("bearer " + auth.AccessToken)
	for _, ce := range lyftEsts.CostEstimates {
		log.Printf("$%0.2f\t%s", float32(ce.EstimatedCostCentsMax/100.0), ce.DisplayName)
	}
}
