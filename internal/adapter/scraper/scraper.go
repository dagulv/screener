package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/dagulv/screener/internal/env"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/rs/xid"
)

type scraper struct {
	browser       *rod.Browser
	env           *env.Environment
	currencyStore port.Currency
	companyStore  port.Company
}

var names = map[string]string{
	"Revenue":                "revenue",
	"Cost of Revenue":        "costOfRevenue",
	"Gross Operating Profit": "grossOperatingProfit",
	"Operating income before interest and taxes": "ebit",
	"Net Income":                        "netIncome",
	"Cash and cash equivalents":         "cashAndEquivalents",
	"Short-term investments":            "shortTermInvestments",
	"Total Assets":                      "totalAssets",
	"Long Term Debt":                    "longTermDebt",
	"Current Debt":                      "currentDebt",
	"Total Liabilities":                 "totalLiabilities",
	"Total stockholders' equity":        "equity",
	"Operating Cash Flow":               "operatingCashFlow",
	"Capital Expenditure":               "capitalExpenditures",
	"Free Cash Flow":                    "freeCashFlow",
	"Not property, plant and equipment": "ppe",
}

func setFieldByJSONTag(data interface{}, jsonTag string, value any) error {
	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")

		// handle cases like `json:"revenue,omitempty"`
		tag = strings.Split(tag, ",")[0]

		if tag == jsonTag {
			fieldValue := v.Field(i)
			if !fieldValue.CanSet() {
				return fmt.Errorf("cannot set field %s", field.Name)
			}

			val := reflect.ValueOf(value)

			// Handle assignability or convertibility
			if val.Type().AssignableTo(fieldValue.Type()) {
				fieldValue.Set(val)
				return nil
			} else if val.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(val.Convert(fieldValue.Type()))
				return nil
			} else {
				return fmt.Errorf("value of type %T not assignable to field %s", value, field.Name)
			}
		}
	}
	return nil // or error if field not found
}

func NewScraper(ctx context.Context, env *env.Environment, currencyStore port.Currency, companyStore port.Company) port.Scraper {
	// c := colly.NewCollector(
	// 	colly.Async(true),
	// )

	// c.OnRequest(func(r *colly.Request) {
	// 	r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	// 	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	// 	r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
	// 	r.Headers.Set("Connection", "keep-alive")
	// 	r.Headers.Set("Upgrade-Insecure-Requests", "1")

	// 	fmt.Println("Visiting", r.URL.String())
	// })

	url := launcher.New().Headless(true).Devtools(true).MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect().Context(ctx)

	// c.OnResponse(func(r *colly.Response) {
	// 	test := r.Body
	// 	_ = test
	// 	fmt.Println("Response", string(r.Body))
	// })

	return scraper{
		browser:       browser,
		env:           env,
		currencyStore: currencyStore,
		companyStore:  companyStore,
	}
}

func (s scraper) GetCompanies(ctx context.Context) (companies []domain.RawCompany, err error) {
	// 1. Read the raw JSON file
	bytes, err := os.ReadFile("response.json")
	if err != nil {
		log.Printf("Failed to read input file: %v", err)
		return
	}

	if err := json.Unmarshal(bytes, &companies); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return nil, err
	}

	return
}

func (s scraper) GetCompanyFinancials(ctx context.Context, companies []domain.Company) iter.Seq2[*domain.Financials, error] {
	return func(yield func(*domain.Financials, error) bool) {
		currencies := make([]domain.IDAndName, 0)

		for c, err := range s.currencyStore.IterateCurrencies(ctx, domain.IDAndNameFilter{}) {
			if err != nil {
				yield(nil, err)
				return
			}

			currencies = append(currencies, *c)
		}

		completedCompanies := make(map[xid.ID]struct{})
		for c, err := range s.companyStore.IterateFinancials(ctx, domain.FinancialFilter{}) {
			if err != nil {
				yield(nil, err)
				return
			}

			completedCompanies[c.CompanyID] = struct{}{}
		}

		if err := rod.Try(func() {
			page := s.browser.MustPage()

			var err error
			err = proto.NetworkEnable{}.Call(page)
			if err != nil {
				panic(err)
			}
			err = proto.NetworkSetCacheDisabled{CacheDisabled: true}.Call(page)
			if err != nil {
				panic(err)
			}

			for {
				completed := true

				for _, c := range companies {
					if _, ok := completedCompanies[c.ID]; ok {
						continue
					}
					completed = false

					var body string
					var exists bool

					wait := page.EachEvent(func(e *proto.NetworkResponseReceived) bool {
						var u *url.URL
						if u, err = url.Parse(e.Response.URL); err != nil {
							return true
						}

						if u.Host == "lt.morningstar.com" && e.Response.MIMEType == "text/html" && e.Response.Status == 200 {
							// Skip first request. For some reason the first is always empty.
							if !exists {
								exists = true
								return false
							}
							// now := time.Now()

							// time.Since(now)
							page.WaitEvent(&proto.NetworkLoadingFinished{})()

							fmt.Printf("Response from %s\n", u.Path)
							var b *proto.NetworkGetResponseBodyResult
							b, err = proto.NetworkGetResponseBody{
								RequestID: e.RequestID,
							}.Call(page)

							fmt.Printf("Error: %s\nStatus Code: %d\n", err, e.Response.Status)
							if err != nil {
								return true
							}

							body = b.Body

							if body == "" {
								err = errors.New("missing body")
								return true
							}
							var doc *goquery.Document
							if doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(body))); err != nil {
								return false
							}

							query := u.Query()
							curr := query.Get("CurrencyId")

							if curr == "" {
								err = errors.New("missing currency")
								return false
							}
							var currencyID xid.ID
							for _, c := range currencies {
								if c.Name == curr {
									currencyID = c.ID
								}
							}

							rows := make([][]string, 0)
							doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
								cols := make([]string, 0)
								s.Find("td,th").Each(func(i int, s *goquery.Selection) {
									cols = append(cols, s.Text())
								})
								rows = append(rows, cols)
							})

							req, _ := http.NewRequest("GET", "https://api.nasdaq.com/api/nordic/instruments/"+c.OrderbookID+"/summary?assetClass=SHARES&lang=en", nil)
							req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
							req.Header.Set("Accept", "application/json")
							req.Header.Set("Referer", "https://nasdaq.com")

							client := &http.Client{}
							var resp *http.Response
							if resp, err = client.Do(req); err != nil {
								return true
							}
							defer resp.Body.Close()

							if resp.StatusCode != http.StatusOK {
								panic(errors.New("bad status: " + resp.Status))
							}

							var result map[string]any
							err := json.NewDecoder(resp.Body).Decode(&result)
							if err != nil {
								log.Fatal(err)
							}

							nosRaw := result["data"].(map[string]any)["summaryData"].(map[string]any)["shares"].(map[string]any)["value"].(string)
							nosFormattedString := strings.ReplaceAll(nosRaw, ",", "")

							var numberOfShares int
							numberOfShares, err = strconv.Atoi(nosFormattedString)

							if err != nil {
								return true
							}

							financials := make([]domain.Financials, 5)

							for i := range financials {
								financials[i].CompanyID = c.ID
								if financials[i].FiscalYear, err = strconv.Atoi(rows[0][i+1]); err != nil {
									return false
								}
								financials[i].CurrencyID = currencyID
								financials[i].StaticData.NumberOfShares = numberOfShares
								filledTags := make(map[string]struct{})
								for _, row := range rows[1:] {
									tag, ok := names[row[0]]

									if !ok {
										continue
									}

									if _, ok := filledTags[tag]; ok {
										continue
									}

									filledTags[tag] = struct{}{}
									var v float64

									if row[i+1] != "" && row[i+1] != "-" {
										formatted := strings.ReplaceAll(row[i+1], ",", "")
										if v, err = strconv.ParseFloat(formatted, 64); err != nil {
											return false
										}
									}

									if err = setFieldByJSONTag(&financials[i].StaticData, tag, int(v*100)); err != nil {
										return false
									}
								}

								if !yield(&financials[i], nil) {
									return true
								}
							}

							return true
						}

						return false
					})
					page.MustWaitIdle()
					page.MustNavigate(s.env.EuropeBaseUrl + "/" + toSlug(c.Symbol) + "/financials")
					fmt.Printf("Navigating %s...\n", s.env.EuropeBaseUrl+"/"+toSlug(c.Symbol)+"/financials")
					wait()

					if err != nil {
						fmt.Println(err)
					}

					completedCompanies[c.ID] = struct{}{}
					completed = true
				}

				if completed {
					break
				}
			}
		}); err != nil {
			yield(nil, err)
			return
		}
		// s.c.OnHTML("nasdaq-market-data-iframe", func(h *colly.HTMLElement) {
		// 	test := h
		// 	_ = test
		// 	if !yield(&domain.Financials{}, nil) {
		// 		return
		// 	}
		// })

		// for _, c := range companies[:1] {
		// 	test := toSlug(c.Symbol)
		// 	_ = test
		// 	if err := s.c.Visit(s.env.EuropeBaseUrl + "/" + toSlug(c.Symbol) + "/financials"); err != nil {
		// 		yield(nil, err)
		// 		return
		// 	}
		// }

		// s.c.Wait()
	}
}

func toSlug(input string) string {
	return strings.ReplaceAll(strings.ToLower(input), " ", "-")
}

func (s scraper) GetCompanyShares(ctx context.Context, companies []domain.Company) iter.Seq2[*domain.Share, error] {
	return func(yield func(*domain.Share, error) bool) {
		completedCompanies := make(map[xid.ID]struct{})
		for c, err := range s.companyStore.IterateShares(ctx, domain.CompanyFilter{}) {
			if err != nil {
				yield(nil, err)
				return
			}

			completedCompanies[c.CompanyID] = struct{}{}
		}

		now := time.Now()
		y, m, d := now.Date()
		pastMonth := time.Date(y, m-1, d, 0, 0, 0, 0, time.UTC)
		for {
			completed := true

			for _, c := range companies {
				if _, ok := completedCompanies[c.ID]; ok {
					continue
				}
				completed = false
				if err := getShare(c.ID, c.OrderbookID, pastMonth, now, yield); err != nil {
					yield(nil, err)
					return
				}
				time.Sleep(time.Millisecond * 500)
				completedCompanies[c.ID] = struct{}{}
				completed = true
			}
			if completed {
				break
			}
		}
	}
}

func (s scraper) GetCompanySharesByFinancials(ctx context.Context, companies []domain.Company) iter.Seq2[*domain.Share, error] {
	return func(yield func(*domain.Share, error) bool) {
		financials := s.companyStore.IterateFinancialsByMissingShare(ctx)

		for f, err := range financials {
			if err != nil {
				yield(nil, err)
				return
			}

			from := time.Date(f.FiscalYear+1, 1, 1, 0, 0, 0, 0, time.UTC)
			to := time.Date(f.FiscalYear+1, 1, 2, 0, 0, 0, 0, time.UTC)
			log.Println(from, to)

			var company domain.Company
			for _, c := range companies {
				if c.ID == f.CompanyID {
					company = c
				}
			}
			if err := getShare(company.ID, company.OrderbookID, from, to, yield); err != nil {
				yield(nil, err)
				return
			}
			time.Sleep(time.Millisecond * 200)
		}
	}
}

func getShare(companyId xid.ID, externalId string, from time.Time, to time.Time, yield func(*domain.Share, error) bool) (err error) {
	req, _ := http.NewRequest("GET", "https://api.nasdaq.com/api/nordic/instruments/"+externalId+"/chart/download?assetClass=SHARES&fromDate="+from.Format(time.DateOnly)+"&toDate="+to.Format(time.DateOnly), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://nasdaq.com")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status: " + resp.Status)
	}

	var data domain.RawRoot
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	for _, raw := range data.Data.Charts.Rows {
		var share domain.Share
		share.CompanyID = companyId
		if share.Date, err = time.Parse(time.DateOnly, raw.DateTime); err != nil {
			return
		}
		if share.Open, err = toFloat(raw.Open); err != nil {
			return
		}
		if share.High, err = toFloat(raw.High); err != nil {
			return
		}
		if share.Low, err = toFloat(raw.Low); err != nil {
			return
		}
		if share.Close, err = toFloat(raw.Close); err != nil {
			return
		}
		if share.Volume, err = toInt(raw.TotalVolume); err != nil {
			return
		}
		if share.Average, err = toFloat(raw.Average); err != nil {
			return
		}

		if !yield(&share, nil) {
			return
		}
	}

	return
}

func (s scraper) GetCompanyMeta(ctx context.Context, company *domain.Company) (err error) {
	time.Sleep(200 * time.Millisecond)
	req, err := http.NewRequest("GET", "https://api.nasdaq.com/api/nordic/instruments/"+company.OrderbookID+"/summary?assetClass=SHARES&lang=en", nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://nasdaq.com")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
		return errors.New("bad status: " + resp.Status)
	}

	var data struct {
		Data struct {
			SummaryData struct {
				InsExecutionVenue struct {
					Value string `json:"value"`
				} `json:"insExecutionVenue"`
			} `json:"summaryData"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	switch strings.ToLower(data.Data.SummaryData.InsExecutionVenue.Value) {
	case string(domain.MarketPlaceCodeXSTO):
		company.MarketPlaceCode = domain.MarketPlaceCodeXSTO
		company.CountryCode = domain.CountryCodeSE
	case string(domain.MarketPlaceCodeXCSE):
		company.MarketPlaceCode = domain.MarketPlaceCodeXCSE
		company.CountryCode = domain.CountryCodeDK
	case string(domain.MarketPlaceCodeXHEL):
		company.MarketPlaceCode = domain.MarketPlaceCodeXHEL
		company.CountryCode = domain.CountryCodeFI
	case string(domain.MarketPlaceCodeXICE):
		company.MarketPlaceCode = domain.MarketPlaceCodeXICE
		company.CountryCode = domain.CountryCodeIS
	}

	return
}

func toFloat(input string) (float64, error) {
	if input == "" {
		return 0, nil
	}

	formatted, err := strconv.ParseFloat(strings.ReplaceAll(input, ",", ""), 64)

	if err != nil {
		return 0, err
	}

	return formatted * 100, nil
}
func toInt(input string) (int, error) {
	if input == "" {
		return 0, nil
	}

	formatted, err := strconv.ParseFloat(strings.ReplaceAll(input, ",", ""), 64)

	if err != nil {
		return 0, err
	}

	return int(formatted * 100), nil
}
