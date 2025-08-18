package domain

import "github.com/rs/xid"

type IDAndName struct {
	ID   xid.ID `json:"id"`
	Name string `json:"name"`
}

type IDAndNameFilter struct {
	Order   string   `query:"order" enum:"asc,desc"`
	OrderBy string   `query:"orderBy"`
	Limit   int      `query:"limit"`
	Offset  int      `query:"order"`
	Include []xid.ID `query:"include"`
	Search  string   `query:"search"`
}
