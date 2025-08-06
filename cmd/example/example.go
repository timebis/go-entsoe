package main

import (
	"fmt"
	"time"

	"github.com/timebis/go-entsoe"
)

func main() {
	client := entsoe.NewEntsoeClientFromEnv()

	dayahead, err := entsoe.NewDayAhead("FR", client, time.Hour)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	prices, err := dayahead.Fetch(time.Now().Add(-7*24*time.Hour), time.Now())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, p := range prices {
		fmt.Printf("%s: %f\n", p.Time.Format("2006-01-02 15:04:05"), p.Price_eur_per_MWh)
	}

}
