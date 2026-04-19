package entsoe

import (
	"strconv"
	"time"
)

const resolution15m = 15 * time.Minute

type DayAhead struct {
	client     *EntsoeClient
	domain     DomainType
	prices     []DayAheadElement
	lastUpdate time.Time
}

type DayAheadElement struct {
	Time              time.Time
	Price_eur_per_MWh float64
}

func NewDayAhead(area Area, client *EntsoeClient) (*DayAhead, error) {
	domain, err := domain(string(area))
	if err != nil {
		return nil, err
	}

	return &DayAhead{
		client:     client,
		domain:     domain,
		lastUpdate: time.Time{},
	}, nil
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
	logger.Info().
		Str("from", from.Format("2006-01-02")).
		Str("to", to.Format("2006-01-02")).
		Msg("Fetching day-ahead prices")

	doc, err := d.client.GetDayAheadPrices(d.domain, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Error fetching day-ahead prices")
		return
	}

	prices, lastUpdate, err := d.parsePublicationMarketDocument(doc)
	if err != nil {
		logger.Error().Err(err).Msg("Error parsing publication market document")
		return
	}

	// merge into d.prices, deduplicating by timestamp
	existing := make(map[int64]struct{}, len(d.prices))
	for _, p := range d.prices {
		existing[p.Time.Unix()] = struct{}{}
	}
	for k, v := range prices {
		if _, dup := existing[k]; dup {
			for _, p := range d.prices {
				if p.Time.Unix() == k && p.Price_eur_per_MWh != v {
					logger.Error().
						Time("slot", time.Unix(k, 0)).
						Float64("existing", p.Price_eur_per_MWh).
						Float64("conflict", v).
						Msg("duplicate slot with different price across fetches")
				}
			}
		} else {
			d.prices = append(d.prices, DayAheadElement{
				Time:              time.Unix(k, 0),
				Price_eur_per_MWh: v,
			})
			existing[k] = struct{}{}
		}
	}

	d.lastUpdate = lastUpdate
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
		step := durations[resolution]

		start, err := time.Parse("2006-01-02T15:04Z", period.TimeInterval.Start)
		if err != nil {
			return nil, time.Time{}, err
		}

		logger.Debug().
			Str("resolution", string(resolution)).
			Str("start", period.TimeInterval.Start).
			Str("end", period.TimeInterval.End).
			Int("points", len(period.Point)).
			Int("slots_per_point", int(step/resolution15m)).
			Msg("time series")

		if step == 0 {
			logger.Warn().Str("resolution", string(resolution)).Msg("unknown resolution, skipping time series")
			continue
		}

		for _, point := range period.Point {
			index, err := strconv.ParseInt(point.Position, 10, 64)
			if err != nil {
				return nil, time.Time{}, err
			}

			price, err := strconv.ParseFloat(point.PriceAmount, 64)
			if err != nil {
				return nil, time.Time{}, err
			}

			pointTime := GetPointTime(start, int(index), resolution)

			// split coarser points into 15-min slots
			slots := int(step / resolution15m)
			for i := 0; i < slots; i++ {
				t := pointTime.Add(time.Duration(i) * resolution15m)
				if existing, exists := res[t.Unix()]; exists {
					if existing != price {
						logger.Warn().
							Time("slot", t).
							Str("area", timeSeries.InDomainMRID.Text).
							Float64("existing", existing).
							Float64("conflict", price).
							Msg("duplicate slot with different price")
					}
				} else {
					res[t.Unix()] = price
				}
			}
		}
	}

	backfill(res)

	return res, end, nil
}

func backfill(res map[int64]float64) {
	if len(res) == 0 {
		return
	}

	var first, last int64
	for t := range res {
		if first == 0 || t < first {
			first = t
		}
		if t > last {
			last = t
		}
	}

	step := int64(resolution15m.Seconds())
	var prevPrice float64

	for t := first; t <= last; t += step {
		if price, ok := res[t]; ok {
			prevPrice = price
		} else {
			logger.Debug().
				Time("slot", time.Unix(t, 0)).
				Float64("backfill_price", prevPrice).
				Msg("missing slot, backfilling with previous price")
			res[t] = prevPrice
		}
	}
}
