package domain

import "github.com/rs/xid"

type MagicRank struct {
	Company           IDAndName `json:"company"`
	FiscalYear        int       `json:"fiscalYear"`
	ROC               float64   `json:"roc"`
	EarningsYield     float64   `json:"earningsYield"`
	ROCRank           int       `json:"rocRank"`
	EarningsYieldRank int       `json:"earningsYieldRank"`
	Rank              int       `json:"rank"`
}

type MagicRankFilter struct {
	Order      string   `query:"order" enum:"asc,desc" default:"asc"`
	OrderBy    string   `query:"orderBy" enum:"rank" default:"rank"`
	Limit      int      `query:"limit" min:"1" max:"500" default:"50"`
	Offset     int      `query:"offset" min:"0"`
	Include    []xid.ID `query:"include"`
	Search     string   `query:"search"`
	FiscalYear int      `query:"fiscalYear"`
}
