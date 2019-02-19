package main

import (
	"fmt"

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
	ubercfg := viper.GetStringMapString("app.uber")
	auth := rideshare.LyftAuth(nil, lyftcfg["username"], lyftcfg["password"])
	lyftEsts := rideshare.LyftCostEstimate("bearer " + auth.AccessToken)
	fmt.Println("--- lyft ---")
	for _, ce := range lyftEsts.CostEstimates {
		fmt.Printf("$%0.2f - $%0.2f\t%s\n", float32(ce.EstimatedCostCentsMin/100.0), float32(ce.EstimatedCostCentsMax/100.0), ce.DisplayName)
	}

	fmt.Println("--- uber ---")
	uberEsts := rideshare.UberCostEstimate("Token " + ubercfg["server_token"])
	for _, ce := range uberEsts.Prices {
		fmt.Printf("$%0.2f - $%0.2f\t%s\n", ce.LowEstimate, ce.HighEstimate, ce.DisplayName)
	}
}
