package route

import (
	"time"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/service"
	"github.com/webmafia/papi"
)

type Screener struct {
	Service service.Screener
}

func (r Screener) IterateScreener(api *papi.API) error {
	type req struct {
		Filter domain.ScreenerFilter
	}

	return papi.GET(api, papi.Route[req, papi.List[domain.Screener]]{
		Path: "/screener",
		Handler: func(ctx *papi.RequestCtx, in *req, out *papi.List[domain.Screener]) (err error) {
			// s, err := session.Current(c)

			// if err != nil {
			// 	return
			// }

			if in.Filter.FiscalYear == 0 {
				now := time.Now()
				y, _, _ := now.Date()
				in.Filter.FiscalYear = y - 1
			}

			count, err := r.Service.CountScreener(ctx, in.Filter)

			if err != nil {
				return
			}

			out.SetTotal(count)

			return out.WriteAll(r.Service.IterateScreener(ctx, in.Filter))
		},
	})
}
