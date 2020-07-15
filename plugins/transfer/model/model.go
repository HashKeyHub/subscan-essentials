package model

import "github.com/shopspring/decimal"

type Transfer struct {
	From           string          `sql:"column:addr_from;default: null;size:100" json:"from"`
	To             string          `sql:"column:addr_to;default: null;size:100" json:"to"`
	BlockNum       int             `json:"block_num" `
	BlockTimestamp int             `json:"block_timestamp"`
	ExtrinsicIndex string          `json:"extrinsic_index" sql:"primary_key;size:100"`
	Fee            decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	Amount         decimal.Decimal `json:"amount" sql:"type:decimal(30,0);"`
	Finalized      bool            `json:"finalized"`
}
