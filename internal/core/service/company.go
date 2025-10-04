package service

import (
	"context"
	"io"
	"iter"
	"strconv"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
	"github.com/xuri/excelize/v2"
)

type Company struct {
	store         port.Company
	currencyStore port.Currency
	sectorStore   port.Sector
	screenerStore port.Screener
}

func NewCompany(store port.Company, currencyStore port.Currency, sectorStore port.Sector, screenerStore port.Screener) Company {
	return Company{
		store:         store,
		currencyStore: currencyStore,
		sectorStore:   sectorStore,
		screenerStore: screenerStore,
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

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (s Company) DownloadFinancials(ctx context.Context, filters domain.ScreenerFilter, w io.Writer) (err error) {
	financials := s.screenerStore.IterateScreener(ctx, filters)

	if err != nil {
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	i := 1
	for financial, err := range financials {
		if err != nil {
			return err
		}
		if i == 1 {
			if err = setTitle(f, financial, i); err != nil {
				return err
			}
			i++
		}
		j := 0
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), financial.Name); err != nil {
			return err
		}
		j++
		if value, ok := financialField(financial.MagicRank); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.Sector); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.CapitalExpenditures); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.CashAndEquivalents); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.CostOfRevenue); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.CurrentDebt); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.EBIT); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.Equity); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.FreeCashFlow); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.GrossOperatingProfit); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.LongTermDebt); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.NetIncome); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.NumberOfShares); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.OperatingCashFlow); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.PPE); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.Revenue); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.ShortTermInvestments); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.TotalAssets); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.TotalLiabilities); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.EPS); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.EVEBIT); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.PB); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.PE); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.PS); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.OperatingMargin); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.NetMargin); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.ROE); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.ROC); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.LiabilitiesToEquity); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.DebtToEbit); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.DebtToAssets); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}
		if value, ok := financialField(financial.CashConversion); ok {
			if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), value); err != nil {
				return err
			}
			j++
		}

		i++
	}

	return f.Write(w)
}

func setTitle(f *excelize.File, financial *domain.Screener, i int) (err error) {
	j := 0
	if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Name"); err != nil {
		return err
	}
	j++
	if _, ok := financialField(financial.MagicRank); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Magic Formula"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.Sector); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Sector"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.CapitalExpenditures); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Capital Expenditures"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.CashAndEquivalents); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Cash and Equivalents"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.CostOfRevenue); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Cost of Revenue"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.CurrentDebt); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Current Debt"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.EBIT); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "EBIT"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.Equity); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Equity"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.FreeCashFlow); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Free Cash Flow"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.GrossOperatingProfit); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Gross Operating Profit"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.LongTermDebt); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Long Term Debt"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.NetIncome); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Net Income"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.NumberOfShares); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Number of Shares"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.OperatingCashFlow); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Operating Cash Flow"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.PPE); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "PPE"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.Revenue); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Revenue"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.ShortTermInvestments); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Short Term Investments"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.TotalAssets); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Total Assets"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.TotalLiabilities); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Total Liabilities"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.EPS); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "EPS"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.EVEBIT); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "EV/EBIT"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.PB); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "P/B"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.PE); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "P/E"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.PS); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "P/S"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.OperatingMargin); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Operating Margin"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.NetMargin); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Net Margin"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.ROE); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "ROE"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.ROC); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "ROC"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.LiabilitiesToEquity); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Liabilities to Equity"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.DebtToEbit); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Debt to Ebit"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.DebtToAssets); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Debt to Assets"); err != nil {
			return err
		}
		j++
	}
	if _, ok := financialField(financial.CashConversion); ok {
		if err = f.SetCellValue("Sheet1", string(alphabet[j])+strconv.Itoa(i), "Cash Conversion Rate"); err != nil {
			return err
		}
		j++
	}
	return
}

func financialField[T any](field domain.Nullable[T]) (T, bool) {
	if field.Valid {
		return field.Content, true
	}

	var t T
	return t, false
}
