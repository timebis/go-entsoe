# go entsoe

Lightweight Go wrapper around the ENTSO-E Transparency Platform RESTful API

## ENTSO-E Transparency Platform

ENTSO-E, the European Network of Transmission System Operators, represents 39 electricity transmission system operators (TSOs) from 35 countries across Europe. The ENTSO-E Transparency Platform aims to provide free, continuous access to pan-European electricity market data for all users, across six main categories: Load, Generation, Transmission, Balancing, Outages and Congestion Management.

## Usage

### Basic usage

Request day ahead prices for the next 24h:

```
client := entsoe.NewEntsoeClient("token")

now := time.Now()
from := now.Truncate(24 * time.Hour)
to := from.AddDate(0, 0, 1)

doc, err := client.GetDayAheadPrices(entsoe.DomainBE, from, to)
```

### Day Ahead Prices

Automatically fetch day-ahead prices at noon every day:

```
dayAhead, err := entsoe.NewDayAhead("BE", "token", time.Hour)

now := time.Now()
price, err := dayAhead.GetDayAheadPrice(now)
```
