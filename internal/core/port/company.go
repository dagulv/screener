package port

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/rs/xid"
)

type Company interface {
	Context

	CreateCompany(ctx context.Context, company *domain.Company) error
	ReadCompany(ctx context.Context, company *domain.Company) error
	UpdateCompany(ctx context.Context, company *domain.Company) error
	DeleteCompany(ctx context.Context, companyId xid.ID) error
	IterateCompanies(ctx context.Context, filters domain.CompanyFilter) iter.Seq2[*domain.Company, error]
	CountCompanies(ctx context.Context, filters domain.CompanyFilter) (int, error)

	CreateFinancials(ctx context.Context, financials *domain.Financials) error
	IterateFinancials(ctx context.Context, filters domain.FinancialFilter) iter.Seq2[*domain.Financials, error]
	CountFinancials(ctx context.Context, filters domain.FinancialFilter) (int, error)
	IterateFinancialsByMissingShare(ctx context.Context) iter.Seq2[*domain.Financials, error]
	CreateShare(ctx context.Context, share *domain.Share) error
	IterateShares(ctx context.Context, filters domain.CompanyFilter) iter.Seq2[*domain.Share, error]
}
