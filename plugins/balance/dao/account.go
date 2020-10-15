package dao

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"math/rand"
)

func ReadStorage(p websocket.WsConn, module, prefix string, hash string, arg ...string) (r storage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	v := &rpc.JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
		return
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return "", nil
		}
		return Decode(dataHex, key.ScaleType, nil)
	}
	return r, err
}

func Decode(raw string, decodeType string, metadata *metadata.Instant) (s storage.StateStorage, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in Decode error is: %v \n", r)
		}
	}()
	m := types.ScaleDecoder{}

	option := types.ScaleDecoderOption{}
	if metadata != nil {
		metadataStruct := types.MetadataStruct(*metadata)
		option.Metadata = &metadataStruct
	}
	option.Spec = 25
	m.Init(types.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
	return storage.StateStorage(util.InterfaceToString(m.ProcessAndUpdateData(decodeType))), nil
}

func getFreeBalance(p websocket.WsConn, accountID, hash string) (decimal.Decimal, decimal.Decimal, error) {
	data, err := ReadStorage(p, "System", "Account", hash, util.TrimHex(accountID))
	if err == nil {
		var account model.AccountData
		log.Info("Account data: %v", data.ToString())
		data.ToAny(&account)
		return account.Data.Free.Add(account.Data.Reserved), decimal.Zero, nil
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
