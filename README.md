# go entsoe

Lightweight Go wrapper around the ENTSO-E Transparency Platform RESTful API

## ENTSO-E Transparency Platform

ENTSO-E, the European Network of Transmission System Operators, represents 39 electricity transmission system operators (TSOs) from 35 countries across Europe. The ENTSO-E Transparency Platform aims to provide free, continuous access to pan-European electricity market data for all users, across six main categories: Load, Generation, Transmission, Balancing, Outages and Congestion Management.

### API Token
One should create an account on the [ENTSO-E Transparency Platform](https://transparency.entsoe.eu/usrm/user/myAccountSettings) and get an API token.  
This token could be stored in the `.env` file or as an environment variable.

## Usage

### Basic usage

Request day ahead prices from the last 7 days:

```go
	client := entsoe.NewEntsoeClientFromEnv()

	dayahead, err := entsoe.NewDayAhead(entsoe.France, client, time.Hour)
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

```