package domain

import (
	"time"

	"github.com/rs/xid"
)

type Share struct {
	CompanyID xid.ID    `json:"companyId"`
	Date      time.Time `json:"date"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int       `json:"volume"`
	Average   float64   `json:"average"`
}

type Row struct {
	DateTime    string `json:"dateTime"`
	Open        string `json:"open"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Close       string `json:"close"`
	TotalVolume string `json:"totalVolume"`
	Average     string `json:"average"`
}

type Charts struct {
	Rows []Row `json:"rows"`
}

type Data struct {
	Charts Charts `json:"charts"`
}

type RawRoot struct {
	Data Data `json:"data"`
}
