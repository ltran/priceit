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
	route := rideshare.Route{
		SLat: 37.7763,
		SLng: -122.3918,
		ELat: 37.7972,
		ELng: -122.4533,
	}

	lyft := rideshare.NewLyft(
		lyftcfg["username"],
		lyftcfg["password"],
	)

	lyftEsts := lyft.GetEstimate(route)
	fmt.Println("--- lyft ---")
	for _, ce := range lyftEsts.CostEstimates {
		fmt.Printf("$%0.2f - $%0.2f\t%s\n", float32(ce.EstimatedCostCentsMin/100.0), float32(ce.EstimatedCostCentsMax/100.0), ce.DisplayName)
	}

	uber := rideshare.NewUber(
		ubercfg["server_token"],
	)

	uberEsts := uber.UberCostEstimate(route)
	fmt.Println("--- uber ---")
	for _, ce := range uberEsts.Prices {
		fmt.Printf("$%0.2f - $%0.2f\t%s\n", ce.LowEstimate, ce.HighEstimate, ce.DisplayName)
	}
}
