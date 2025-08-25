package service

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
)

type Company struct {
	store         port.Company
	currencyStore port.Currency
	sectorStore   port.Sector
}

func NewCompany(store port.Company, currencyStore port.Currency, sectorStore port.Sector) Company {
	return Company{
		store:         store,
		currencyStore: currencyStore,
		sectorStore:   sectorStore,
	}
}

func (s Company) Create(ctx context.Context, company *domain.Company) (err error) {
	company.ID = xid.New()

	return s.store.CreateCompany(ctx, company)
}

func (s Company) Read(ctx context.Context, company *domain.Company) (err error) {
	return s.store.ReadCompany(ctx, company)
}

func (s Company) Update(ctx context.Context, company *domain.Company) (err error) {
	return s.store.UpdateCompany(ctx, company)
}

func (s Company) Delete(ctx context.Context, companyId xid.ID) (err error) {
	return s.store.DeleteCompany(ctx, companyId)
}

func (s Company) Count(ctx context.Context, filters domain.CompanyFilter) (int, error) {
	return s.store.CountCompanies(ctx, filters)
}

func (s Company) Iterate(ctx context.Context, filters domain.CompanyFilter) iter.Seq2[*domain.Company, error] {
	return s.store.IterateCompanies(ctx, filters)
}

// func (s Company) ImportCompanies(ctx context.Context) (err error) {
// 	rawCompanies, err := s.scraper.GetCompanies(ctx)

// 	if ctx, err = s.store.AcquireContext(ctx); err != nil {
// 		return
// 	}
// 	defer s.store.ReleaseContext(ctx)

// 	companies := make([]domain.Company, 0)
// 	currencies := make([]domain.IDAndName, 0)
// 	sectors := make([]domain.IDAndName, 0)
// 	now := time.Now()

// 	for c, err := range s.currencyStore.IterateCurrencies(ctx, domain.IDAndNameFilter{}) {
// 		if err != nil {
// 			return err
// 		}

// 		currencies = append(currencies, *c)
// 	}
// 	for s, err := range s.sectorStore.IterateSectors(ctx, domain.IDAndNameFilter{}) {
// 		if err != nil {
// 			return err
// 		}

// 		sectors = append(sectors, *s)
// 	}

// 	for c, err := range s.store.IterateCompanies(ctx, domain.CompanyFilter{}) {
// 		if err != nil {
// 			return err
// 		}

// 		companies = append(companies, *c)
// 	}

// 	for _, rawCompany := range rawCompanies {
// 		if slices.ContainsFunc(companies, func(c domain.Company) bool {
// 			return c.ISIN == rawCompany.ISIN
// 		}) {
// 			continue
// 		}

// 		company := domain.Company{
// 			Name: rawCompany.Name,
// 			Bio: domain.Nullable[string]{
// 				Content: rawCompany.Bio,
// 				Valid:   rawCompany.Bio != "",
// 			},
// 			Symbol:      rawCompany.Symbol,
// 			ISIN:        rawCompany.ISIN,
// 			OrderbookID: rawCompany.OrderbookID,
// 		}

// 		for _, currency := range currencies {
// 			if currency.Name == rawCompany.Currency {
// 				company.Currency.ID = currency.ID
// 				break
// 			}
// 		}

// 		if company.Currency.ID.IsNil() {
// 			currency := domain.IDAndName{
// 				ID:   xid.NewWithTime(now),
// 				Name: rawCompany.Currency,
// 			}
// 			if err = s.currencyStore.CreateCurrency(ctx, &currency); err != nil {
// 				return err
// 			}
// 			company.Currency.ID = currency.ID
// 			currencies = append(currencies, currency)
// 		}

// 		for _, sector := range sectors {
// 			if sector.Name == rawCompany.Sector {
// 				company.Sector.ID = sector.ID
// 				break
// 			}
// 		}

// 		if company.Sector.ID.IsNil() {
// 			sector := domain.IDAndName{
// 				ID:   xid.NewWithTime(now),
// 				Name: rawCompany.Sector,
// 			}
// 			if err = s.sectorStore.CreateSector(ctx, &sector); err != nil {
// 				return err
// 			}
// 			company.Sector.ID = sector.ID
// 			sectors = append(sectors, sector)
// 		}

// 		if err = s.Create(ctx, &company); err != nil {
// 			return err
// 		}
// 	}

// 	return s.store.CommitContext(ctx)
// }

// func (s Company) ImportCompanyFinancials(ctx context.Context) (err error) {
// 	companies := make([]domain.Company, 0)
// 	for c, err := range s.store.IterateCompanies(ctx, domain.CompanyFilter{}) {
// 		if err != nil {
// 			return err
// 		}

// 		companies = append(companies, *c)
// 	}

// 	financials := s.scraper.GetCompanyFinancials(ctx, companies)

// 	// if ctx, err = s.store.AcquireContext(ctx); err != nil {
// 	// 	return
// 	// }
// 	// defer s.store.ReleaseContext(ctx)

// 	for f, err := range financials {
// 		if err != nil {
// 			return err
// 		}

// 		if err = s.store.CreateFinancials(ctx, f); err != nil {
// 			return err
// 		}
// 	}
// 	return
// 	// return s.store.CommitContext(ctx)
// }

// func (s Company) ImportCompanyShares(ctx context.Context) (err error) {
// 	companies := make([]domain.Company, 0)
// 	for c, err := range s.store.IterateCompanies(ctx, domain.CompanyFilter{}) {
// 		if err != nil {
// 			return err
// 		}

// 		companies = append(companies, *c)
// 	}

// 	shares := s.scraper.GetCompanyShares(ctx, companies)

// 	for sh, err := range shares {
// 		if err != nil {
// 			return err
// 		}

// 		if err = s.store.CreateShare(ctx, sh); err != nil {
// 			return err
// 		}
// 	}

// 	return
// }

// func (s Company) ImportCompanySharesByFinancials(ctx context.Context) (err error) {
// 	companies := make([]domain.Company, 0)
// 	for c, err := range s.store.IterateCompanies(ctx, domain.CompanyFilter{}) {
// 		if err != nil {
// 			return err
// 		}

// 		companies = append(companies, *c)
// 	}

// 	shares := s.scraper.GetCompanySharesByFinancials(ctx, companies)

// 	for sh, err := range shares {
// 		if err != nil {
// 			return err
// 		}

// 		if err = s.store.CreateShare(ctx, sh); err != nil {
// 			return err
// 		}
// 		log.Println("Created share", sh.Date, sh.CompanyID)
// 	}

// 	return
// }

// func (s Company) ImportCompanyMeta(ctx context.Context) (err error) {
// 	for c, err := range s.store.IterateCompanies(ctx, domain.CompanyFilter{OrderBy: "name", Limit: 10000000}) {
// 		if err != nil {
// 			return err
// 		}

// 		if c.CountryCode != "" && c.MarketPlaceCode != "" {
// 			continue
// 		}

// 		if err = s.scraper.GetCompanyMeta(ctx, c); err != nil {
// 			return err
// 		}

// 		if err = s.store.UpdateCompany(ctx, c); err != nil {
// 			return err
// 		}
// 	}

// 	return
// }

func (s Company) CountFinancials(ctx context.Context, filters domain.FinancialFilter) (int, error) {
	return s.store.CountFinancials(ctx, filters)
}

func (s Company) IterateFinancials(ctx context.Context, filters domain.FinancialFilter) iter.Seq2[*domain.Financials, error] {
	return s.store.IterateFinancials(ctx, filters)
}
