package dao

import (
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
)

type AccountData struct {
	Free         decimal.Decimal `json:"free"`
	Reserved     decimal.Decimal `json:"reserved"`
	FreeKton     decimal.Decimal `json:"free_kton,omitempty"`
	ReservedKton decimal.Decimal `json:"reserved_kton,omitempty"`
	MiscFrozen   decimal.Decimal `json:"misc_frozen"`
	FeeFrozen    decimal.Decimal `json:"fee_frozen"`
}

func getFreeBalance(p websocket.WsConn, accountID, hash string) (decimal.Decimal, decimal.Decimal, error) {
	data, err := rpc.ReadStorage(p, "System", "Account", hash, util.TrimHex(accountID))
	if err == nil {
		var account AccountData
		data.ToAny(&account)
		return account.Free.Add(account.Reserved), decimal.Zero, nil
	}

	return decimal.Zero, decimal.Zero, err
}

func GetBalanceFromNetwork(address string) (decimal.Decimal, error) {
	balance, _, err := getFreeBalance(nil, address, "")
	if err != nil {
		log.Error("GetBalanceFromNetwork error %v", err)
		return decimal.Zero, err
	}
	return balance, nil
}
