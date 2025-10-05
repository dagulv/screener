package domain

const BaseCurrency = "USD"

type CurrencyRate struct {
	FiscalYear int
	Quarter    int
	Currency   IDAndName
	Rate       float32
}

type CurrencyRateResponse struct {
	Rates map[string]map[string]float32 `json:"rates"`
}
