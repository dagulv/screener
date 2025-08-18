package domain

import "github.com/rs/xid"

type RawCompany struct {
	Name        string `json:"fullName"`
	Bio         string `json:"bio"`
	Symbol      string `json:"symbol"`
	ISIN        string `json:"isin"`
	Currency    string `json:"currency"`
	Sector      string `json:"sector"`
	OrderbookID string `json:"orderbookId"`
}

type Company struct {
	ID              xid.ID           `json:"id"`
	Name            string           `json:"name"`
	Bio             Nullable[string] `json:"bio"`
	Symbol          string           `json:"symbol"`
	ISIN            string           `json:"isin"`
	Currency        IDAndName        `json:"currency"`
	Sector          IDAndName        `json:"sector"`
	OrderbookID     string           `json:"orderbookId"`
	CountryCode     CountryCode      `json:"countryCode"`
	MarketPlaceCode MarketPlaceCode  `json:"marketPlaceCode"`
}

type CompanyFilter struct {
	Order   string   `query:"order" enum:"asc,desc" default:"asc"`
	OrderBy string   `query:"orderBy" enum:"name" default:"name"`
	Limit   int      `query:"limit" min:"1" max:"500" default:"50"`
	Offset  int      `query:"offset" min:"0"`
	Include []xid.ID `query:"include"`
	Search  string   `query:"search"`
}

type CountryCode string

const (
	CountryCodeSE CountryCode = "se"
	CountryCodeDK CountryCode = "dk"
	CountryCodeFI CountryCode = "fi"
	CountryCodeIS CountryCode = "is"
)

type MarketPlaceCode string

const (
	MarketPlaceCodeXSTO MarketPlaceCode = "xsto"
	MarketPlaceCodeXCSE MarketPlaceCode = "xcse"
	MarketPlaceCodeXHEL MarketPlaceCode = "xhel"
	MarketPlaceCodeXICE MarketPlaceCode = "xice"
)
