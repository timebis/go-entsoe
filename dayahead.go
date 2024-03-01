package entsoe

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type DayAhead struct {
	client     *EntsoeClient
	domain     DomainType
	resolution time.Duration
	prices     map[int64]float64
	lastUpdate time.Time
	mu         sync.RWMutex
}

func NewDayAhead(area, token string, resolution time.Duration) (*DayAhead, error) {
	if token == "" {
		return nil, fmt.Errorf("missing token")
	}

	domain, err := domain(area)
	if err != nil {
		return nil, err
	}

	if resolution != time.Hour && resolution != 30*time.Minute && resolution != 15*time.Minute {
		return nil, fmt.Errorf("unsupported resolution %s", resolution)
	}

	d := &DayAhead{
		client:     NewEntsoeClient(token),
		domain:     domain,
		resolution: resolution,
		lastUpdate: time.Time{},
	}

	go d.run()

	return d, nil
}

func (d *DayAhead) GetDayAheadPrice(u time.Time) (float64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	price, ok := d.prices[u.Truncate(d.resolution).Unix()]
	if !ok {
		return 0, fmt.Errorf("no price point found for time %s", u)
	}

	return price, nil
}

func (d *DayAhead) run() {
	for {
		now := time.Now()
		from := now.Truncate(24 * time.Hour)
		to := from.AddDate(0, 0, 1)

		doc, err := d.client.GetDayAheadPrices(d.domain, from, to)
		if err == nil {
			prices, lastUpdate, err := d.parsePublicationMarketDocument(doc)
			if err == nil {
				d.mu.Lock()
				d.prices = prices
				d.lastUpdate = lastUpdate
				d.mu.Unlock()
			}
		}

		timer := time.NewTimer(d.nextUpdate())
		<-timer.C
	}
}

func (d *DayAhead) nextUpdate() time.Duration {
	// Day-ahead prices for the next day are available at noon UTC
	next := time.Date(d.lastUpdate.Year(), d.lastUpdate.Month(), d.lastUpdate.Day(), 12, 5, 0, 0, time.UTC)

	now := time.Now().UTC()
	if now.After(next) {
		// New data should already be available, retry in 15 minutes
		return 15 * time.Minute
	}

	return next.Sub(now)
}

func (d *DayAhead) parsePublicationMarketDocument(doc *PublicationMarketDocument) (map[int64]float64, time.Time, error) {
	end, err := time.Parse("2006-01-02T15:04Z", doc.PeriodTimeInterval.End)
	if err != nil {
		return nil, time.Time{}, err
	}

	res := make(map[int64]float64)

	for _, timeSeries := range doc.TimeSeries {
		period := timeSeries.Period
		resolution := ResolutionType(period.Resolution)

		if durations[resolution] != d.resolution {
			// Skip other resolutions
			continue
		}

		start, err := time.Parse("2006-01-02T15:04Z", period.TimeInterval.Start)
		if err != nil {
			return nil, time.Time{}, err
		}

		points := period.Point
		for _, point := range points {
			index, err := strconv.ParseInt(point.Position, 10, 64)
			if err != nil {
				return nil, time.Time{}, err
			}

			price, err := strconv.ParseFloat(point.PriceAmount, 64)
			if err != nil {
				return nil, time.Time{}, err
			}

			res[GetPointTime(start, int(index), resolution).Unix()] = price
		}
	}

	return res, end, nil
}
