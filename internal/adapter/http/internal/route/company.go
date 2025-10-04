package route

import (
	"time"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/service"
	"github.com/rs/xid"
	"github.com/webmafia/papi"
)

type Company struct {
	Service service.Company
}

func (r Company) CreateCompany(api *papi.API) error {
	type req struct {
		Body domain.Company `body:"json"`
	}

	return papi.POST(api, papi.Route[req, domain.Company]{
		Path: "/companies",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			err = r.Service.Create(ctx, &in.Body)
			*out = in.Body
			return
		},
	})
}

func (r Company) GetCompany(api *papi.API) error {
	type req struct {
		CompanyID xid.ID `param:"id"`
	}

	return papi.GET(api, papi.Route[req, domain.Company]{
		Path: "/companies/{id}",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			out.ID = in.CompanyID
			return r.Service.Read(ctx, out)
		},
	})
}

func (r Company) UpdateCompany(api *papi.API) error {
	type req struct {
		CompanyID xid.ID         `param:"id"`
		Body      domain.Company `body:"json"`
	}

	return papi.PUT(api, papi.Route[req, domain.Company]{
		Path: "/companies/{id}",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			// s, err := session.Current(ctx)

			// if err != nil {
			// 	return
			// }

			in.Body.ID = in.CompanyID
			err = r.Service.Update(ctx, &in.Body)
			*out = in.Body
			return
		},
	})
}

func (r Company) DeleteCompany(api *papi.API) error {
	type req struct {
		CompanyId xid.ID
	}

	return papi.DELETE(api, papi.Route[req, domain.Company]{
		Path: "/companies/{id}",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			return r.Service.Delete(ctx, in.CompanyId)
		},
	})
}

func (r Company) IterateCompanies(api *papi.API) error {
	type req struct {
		Filter domain.CompanyFilter
	}

	return papi.GET(api, papi.Route[req, papi.List[domain.Company]]{
		Path: "/companies",
		Handler: func(ctx *papi.RequestCtx, in *req, out *papi.List[domain.Company]) (err error) {
			// s, err := session.Current(c)

			// if err != nil {
			// 	return
			// }

			count, err := r.Service.Count(ctx, in.Filter)

			if err != nil {
				return
			}

			out.SetTotal(count)

			return out.WriteAll(r.Service.Iterate(ctx, in.Filter))
		},
	})
}

func (r Company) IterateFinancials(api *papi.API) error {
	type req struct {
		CompanyId xid.ID `param:"id"`
		Filter    domain.FinancialFilter
	}

	return papi.GET(api, papi.Route[req, papi.List[domain.Financials]]{
		Path: "/financials/{id}",
		Handler: func(ctx *papi.RequestCtx, in *req, out *papi.List[domain.Financials]) (err error) {
			// s, err := session.Current(c)

			// if err != nil {
			// 	return
			// }

			in.Filter.Include = []xid.ID{in.CompanyId}

			count, err := r.Service.CountFinancials(ctx, in.Filter)

			if err != nil {
				return
			}

			out.SetTotal(count)

			return out.WriteAll(r.Service.IterateFinancials(ctx, in.Filter))
		},
	})
}

func (r Company) DownloadFinancials(api *papi.API) (err error) {
	type req struct {
		Filter domain.ScreenerFilter
	}

	return papi.GET(api, papi.Route[req, papi.File[domain.FinancialsFile]]{
		Path: "/financials/download",

		Handler: func(ctx *papi.RequestCtx, in *req, out *papi.File[domain.FinancialsFile]) (err error) {
			if in.Filter.FiscalYear == 0 {
				now := time.Now()
				y, _, _ := now.Date()
				in.Filter.FiscalYear = y - 1
			}

			in.Filter.OrderBy = "name"
			in.Filter.Limit = 1000

			out.SetFilename("financials.xlsx")

			return r.Service.DownloadFinancials(ctx, in.Filter, out.Writer())
		},
	})
}
