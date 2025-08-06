package entsoe

import (
	"fmt"
	"strconv"
	"time"
)

type DayAhead struct {
	client     *EntsoeClient
	domain     DomainType
	resolution time.Duration
	prices     []DayAheadElement
	lastUpdate time.Time
}

type DayAheadElement struct {
	Time              time.Time
	Price_eur_per_MWh float64
}

func NewDayAhead(area string, client *EntsoeClient, resolution time.Duration) (*DayAhead, error) {
	domain, err := domain(area)
	if err != nil {
		return nil, err
	}

	if resolution != time.Hour && resolution != 30*time.Minute && resolution != 15*time.Minute {
		return nil, fmt.Errorf("unsupported resolution %s", resolution)
	}

	d := &DayAhead{
		client:     client,
		domain:     domain,
		resolution: resolution,
		lastUpdate: time.Time{},
	}

	return d, nil
}

func (d *DayAhead) Fetch(from, to time.Time) ([]DayAheadElement, error) {

	for to.Sub(from) > 30*24*time.Hour {
		fromChunk := to.Add(-30 * 24 * time.Hour)
		d.fetch(fromChunk, to)
		to = fromChunk
	}

	d.fetch(from, to)

	return d.prices, nil
}

func (d *DayAhead) fetch(from, to time.Time) {
	fmt.Printf("Fetching from %s to %s\n", from.Format("2006-01-02"), to.Format("2006-01-02"))
	doc, err := d.client.GetDayAheadPrices(d.domain, from, to)
	if err == nil {
		prices, lastUpdate, err := d.parsePublicationMarketDocument(doc)
		if err != nil {
			fmt.Printf("Error parsing publication market document: %s\n", err)
			return
		}

		for k, v := range prices {
			d.prices = append(d.prices, DayAheadElement{
				Time:              time.Unix(k, 0),
				Price_eur_per_MWh: v,
			})
		}

		d.lastUpdate = lastUpdate
	}
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
