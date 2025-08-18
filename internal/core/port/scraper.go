package port

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
)

type Scraper interface {
	GetCompanies(ctx context.Context) ([]domain.RawCompany, error)
	GetCompanyFinancials(ctx context.Context, companies []domain.Company) iter.Seq2[*domain.Financials, error]
	GetCompanyShares(ctx context.Context, companies []domain.Company) iter.Seq2[*domain.Share, error]
	GetCompanySharesByFinancials(ctx context.Context, companies []domain.Company) iter.Seq2[*domain.Share, error]
	GetCompanyMeta(ctx context.Context, company *domain.Company) error
}
