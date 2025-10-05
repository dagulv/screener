package main

import (
	"context"
	"log"
	"sync"

	"github.com/dagulv/screener/internal/adapter/cron"
	"github.com/dagulv/screener/internal/adapter/http"
	"github.com/dagulv/screener/internal/adapter/postgres"
	"github.com/dagulv/screener/internal/core/service"
	"github.com/dagulv/screener/internal/env"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := start(ctx); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) (err error) {
	env, err := env.Load(true)

	if err != nil {
		return
	}

	db, err := postgres.NewDB(ctx, env)

	if err != nil {
		return
	}

	currencyStore := postgres.NewCurrency(db)
	sectorStore := postgres.NewSector(db)
	companyStore := postgres.NewCompany(db)
	// scraper := scraper.NewScraper(ctx, env, currencyStore, companyStore)
	screenerStore := postgres.NewScreener(db)

	scheduler, err := cron.New()

	if err != nil {
		return
	}

	currencyService := service.NewCurrency(currencyStore, env, scheduler)

	if err = currencyService.StartJobs(ctx); err != nil {
		return
	}

	service := http.Service{
		Company:  service.NewCompany(companyStore, currencyStore, sectorStore, screenerStore),
		Screener: service.NewScreener(screenerStore),
	}

	api, err := http.NewApi(env, service, nil)

	if err != nil {
		return
	}

	if err = api.ApiDocsToFile("openapi.json"); err != nil {
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := api.Start(); err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	return api.Close()
}
