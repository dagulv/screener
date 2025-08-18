package route

import (
	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/service"
	"github.com/rs/xid"
	"github.com/webmafia/papi"
)

type Company struct {
	Service service.Company
}

func (s Company) CreateCompany(api *papi.API) error {
	type req struct {
		Body domain.Company `body:"json"`
	}

	return papi.POST(api, papi.Route[req, domain.Company]{
		Path: "/companies",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			err = s.Service.Create(ctx, &in.Body)
			*out = in.Body
			return
		},
	})
}

func (s Company) GetCompany(api *papi.API) error {
	type req struct {
		CompanyID xid.ID `param:"id"`
	}

	return papi.GET(api, papi.Route[req, domain.Company]{
		Path: "/companies/{id}",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			out.ID = in.CompanyID
			return s.Service.Read(ctx, out)
		},
	})
}

func (s Company) UpdateCompany(api *papi.API) error {
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
			err = s.Service.Update(ctx, &in.Body)
			*out = in.Body
			return
		},
	})
}

func (s Company) DeleteCompany(api *papi.API) error {
	type req struct {
		CompanyId xid.ID
	}

	return papi.DELETE(api, papi.Route[req, domain.Company]{
		Path: "/companies/{id}",
		Handler: func(ctx *papi.RequestCtx, in *req, out *domain.Company) (err error) {
			return s.Service.Delete(ctx, in.CompanyId)
		},
	})
}

func (s Company) IterateCompanies(api *papi.API) error {
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

			count, err := s.Service.Count(ctx, in.Filter)

			if err != nil {
				return
			}

			out.SetTotal(count)

			return out.WriteAll(s.Service.Iterate(ctx, in.Filter))
		},
	})
}

func (s Company) IterateFinancials(api *papi.API) error {
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

			count, err := s.Service.CountFinancials(ctx, in.Filter)

			if err != nil {
				return
			}

			out.SetTotal(count)

			return out.WriteAll(s.Service.IterateFinancials(ctx, in.Filter))
		},
	})
}
