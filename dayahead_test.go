package entsoe_test

// not real tests, just examples.

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/timebis/go-entsoe"
)

var log = zerolog.New(
	zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	},
).With().Timestamp().Logger()

func TestDayAheadFetchParis(t *testing.T) {
	client := entsoe.NewEntsoeClientFromEnv()

	da, err := entsoe.NewDayAhead(entsoe.France, client)
	if err != nil {
		t.Fatal(err)
	}

	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatal(err)
	}

	from := time.Date(2026, 3, 1, 0, 0, 0, 0, paris)
	to := time.Date(2026, 3, 2, 0, 0, 0, 0, paris)

	prices, err := da.Fetch(from, to)
	if err != nil {
		t.Fatal(err)
	}

	log.Info().Int("count", len(prices)).Msg("fetched")
	if len(prices) == 0 {
		t.Fatal("no prices returned")
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Time.Before(prices[j].Time)
	})

	for _, p := range prices {
		log.Info().
			Str("slot", p.Time.In(paris).Format("2006-01-02 15:04")).
			Float64("price_eur_per_mwh", p.Price_eur_per_MWh).
			Msg("price")
	}
}

func TestDayAheadFetchLogPoints(t *testing.T) {
	client := entsoe.NewEntsoeClientFromEnv()

	da, err := entsoe.NewDayAhead(entsoe.France, client)
	if err != nil {
		t.Fatal(err)
	}

	from := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)

	prices, err := da.Fetch(from, to)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range prices {
		log.Info().
			Str("slot", p.Time.Format("2006-01-02 15:04")).
			Float64("price_eur_per_mwh", p.Price_eur_per_MWh).
			Msg("price")
	}
}

func TestDayAheadFetchLast7Days(t *testing.T) {
	client := entsoe.NewEntsoeClientFromEnv()

	da, err := entsoe.NewDayAhead(entsoe.Germany, client)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	prices, err := da.Fetch(now.Add(-7*24*time.Hour), now)
	if err != nil {
		t.Fatal(err)
	}

	log.Info().Int("count", len(prices)).Str("area", "Germany").Msg("fetched prices")
}

func TestDayAheadFetchHalfHourResolution(t *testing.T) {
	client := entsoe.NewEntsoeClientFromEnv()

	da, err := entsoe.NewDayAhead(entsoe.Belgium, client)
	if err != nil {
		t.Fatal(err)
	}

	from := time.Date(2025, 4, 18, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 4, 19, 0, 0, 0, 0, time.UTC)

	prices, err := da.Fetch(from, to)
	if err != nil {
		t.Fatal(err)
	}

	log.Info().Int("count", len(prices)).Str("area", "Belgium").Msg("fetched prices")
}

func TestDayAheadBackfill(t *testing.T) {
	client := entsoe.NewEntsoeClientFromEnv()

	da, err := entsoe.NewDayAhead(entsoe.France, client)
	if err != nil {
		t.Fatal(err)
	}

	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Fatal(err)
	}

	from := time.Date(2026, 3, 1, 0, 0, 0, 0, paris)
	to := time.Date(2026, 3, 2, 0, 0, 0, 0, paris)

	prices, err := da.Fetch(from, to)
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Time.Before(prices[j].Time)
	})

	// verify no gaps: every consecutive pair must be exactly 15 min apart
	for i := 1; i < len(prices); i++ {
		gap := prices[i].Time.Sub(prices[i-1].Time)
		if gap != 15*time.Minute {
			t.Errorf("gap detected between %s and %s (%s)",
				prices[i-1].Time.In(paris).Format("2006-01-02 15:04"),
				prices[i].Time.In(paris).Format("2006-01-02 15:04"),
				gap,
			)
		}
	}

	log.Info().Int("count", len(prices)).Msg("no gaps found after backfill")
}

func TestDayAheadInvalidArea(t *testing.T) {
	client := entsoe.NewEntsoeClientFromEnv()

	_, err := entsoe.NewDayAhead("XX", client)
	if err == nil {
		t.Fatal("expected error for unsupported area")
	}
	log.Info().Err(err).Msg("got expected error")
}
