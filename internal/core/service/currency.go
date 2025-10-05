package service

import (
	"context"
	"encoding/json"
	"iter"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/dagulv/screener/internal/env"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/xid"
)

type Currency struct {
	store     port.Currency
	env       *env.Environment
	scheduler gocron.Scheduler
}

func NewCurrency(store port.Currency, env *env.Environment, scheduler gocron.Scheduler) Currency {
	return Currency{
		store:     store,
		env:       env,
		scheduler: scheduler,
	}
}

func (s Currency) StartJobs(ctx context.Context) (err error) {
	_, err = s.scheduler.NewJob(gocron.MonthlyJob(3, gocron.NewDaysOfTheMonth(-1), gocron.NewAtTimes(gocron.NewAtTime(23, 0, 0))), gocron.NewTask(s.FetchCurrencyRates), gocron.WithContext(ctx))
	if err != nil {
		return
	}
	s.scheduler.Start()
	// if err = job.RunNow(); err != nil {
	// 	return
	// }
	return
}

func (s Currency) Create(ctx context.Context, currency *domain.IDAndName) (err error) {
	currency.ID = xid.New()

	return s.store.CreateCurrency(ctx, currency)
}

func (s Currency) Read(ctx context.Context, currency *domain.IDAndName) (err error) {
	return s.store.ReadCurrency(ctx, currency)
}

func (s Currency) Update(ctx context.Context, currency *domain.IDAndName) (err error) {
	return s.store.UpdateCurrency(ctx, currency)
}

func (s Currency) Delete(ctx context.Context, currencyId xid.ID) (err error) {
	return s.store.DeleteCurrency(ctx, currencyId)
}

func (s Currency) Iterate(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error] {
	return s.store.IterateCurrencies(ctx, filters)
}

type sumCount struct {
	sum   float32
	count int
}

func (s Currency) FetchCurrencyRates(ctx context.Context) (err error) {
	currencies := make([]*domain.IDAndName, 0)
	currencySymbols := make([]string, 0)
	oldRates := make([]*domain.CurrencyRate, 0)
	const startYear = 2020
	now := time.Now()

	for c, err := range s.store.IterateCurrencies(ctx, domain.IDAndNameFilter{}) {
		if err != nil {
			return err
		}
		currencies = append(currencies, c)
		currencySymbols = append(currencySymbols, c.Name)
	}

	for r, err := range s.store.IterateCurrencyRates(ctx, domain.IDAndNameFilter{}) {
		if err != nil {
			return err
		}
		oldRates = append(oldRates, r)
	}

	for i := range now.Year() - startYear {
		year := startYear + i
		skip := true

		for _, c := range currencies {
			var quartersSum int
			for _, r := range oldRates {
				if r.FiscalYear != year || r.Currency.ID != c.ID {
					continue
				}

				if quartersSum == 10 {
					break
				}

				quartersSum += r.Quarter
			}

			if quartersSum < 10 {
				skip = false
				break
			}
		}

		if skip {
			continue
		}

		endQuarters := [4]time.Time{
			time.Date(year, time.March, 31, 23, 0, 0, 0, time.UTC),
			time.Date(year, time.June, 30, 23, 0, 0, 0, time.UTC),
			time.Date(year, time.September, 30, 23, 0, 0, 0, time.UTC),
			time.Date(year, time.December, 31, 23, 0, 0, 0, time.UTC),
		}

		start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(year, time.December, 31, 23, 59, 0, 0, time.UTC)

		if end.After(now) {
			for _, q := range endQuarters {
				if now.After(q) {
					end = q
					break
				}
			}
		}

		url, err := url.Parse(s.env.CurrencyRateEndpoint + "/" + start.Format(time.DateOnly) + ".." + end.Format(time.DateOnly))

		if err != nil {
			return err
		}

		q := url.Query()
		q.Set("base", domain.BaseCurrency)
		q.Set("symbols", strings.Join(currencySymbols, ","))
		url.RawQuery = q.Encode()

		log.Println("fetching " + url.String())
		resp, err := http.Get(url.String())

		if err != nil {
			return err
		}

		var result domain.CurrencyRateResponse
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		newRates := make([]domain.CurrencyRate, 0, len(currencies)*4)
		sumNewRates := make(map[string]*[4]sumCount)

		for _, c := range currencies {
			for q := range 4 {
				cr := domain.CurrencyRate{
					FiscalYear: year,
					Quarter:    q,
					Currency:   *c,
				}
				newRates = append(newRates, cr)
			}
			sumNewRates[c.Name] = &[4]sumCount{}
		}

		for ts, r := range result.Rates {
			t, err := time.Parse(time.DateOnly, ts)
			if err != nil {
				return err
			}
			for s, rate := range r {
				_, ok := sumNewRates[s]

				if !ok {
					continue
				}

				for i, q := range endQuarters {
					if t.After(q) {
						continue
					}

					sumNewRates[s][i].count++
					sumNewRates[s][i].sum += rate
					break
				}
			}
		}

		for symbol, sum := range sumNewRates {
			for i := range newRates {
				if newRates[i].Currency.Name != symbol {
					continue
				}

				newRates[i].Rate = sum[newRates[i].Quarter].sum / float32(sum[newRates[i].Quarter].count)
			}
		}

		if err = s.store.SetCurrencyRates(ctx, newRates); err != nil {
			return err
		}
	}

	return
}
