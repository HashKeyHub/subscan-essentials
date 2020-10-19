package dao

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/go-kratos/kratos/pkg/log"
	ws "github.com/gorilla/websocket"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
)

func ReadStorage(p websocket.WsConn, module, prefix string, hash string, spec int, arg ...string) (r storage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	v := &rpc.JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
		log.Warn("got error: (type: %T) %v, try to resend message", err, err)
		if ws.IsUnexpectedCloseError(err, ws.CloseInvalidFramePayloadData,
			ws.CloseMessageTooBig, ws.CloseProtocolError) != true {
			// resend
			if err = websocket.SendWsRequest(p, v, rpc.StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
				log.Error("send ws failed: %v", err)
				return
			}
		} else {
			log.Error("send ws failed: %v", err)
			return
		}
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return "", nil
		}
		return Decode(dataHex, key.ScaleType, nil, spec)
	}
	return r, err
}

func Decode(raw string, decodeType string, metadata *metadata.Instant, spec int) (s storage.StateStorage, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovering from panic in Decode error is: %v", r)
		}
	}()
	m := types.ScaleDecoder{}

	option := types.ScaleDecoderOption{}
	if metadata != nil {
		metadataStruct := types.MetadataStruct(*metadata)
		option.Metadata = &metadataStruct
	}
	option.Spec = spec
	m.Init(types.ScaleBytes{Data: util.HexToBytes(raw)}, &option)
	return storage.StateStorage(util.InterfaceToString(m.ProcessAndUpdateData(decodeType))), nil
}

func getFreeBalance(p websocket.WsConn, accountID, hash string, spec int) (decimal.Decimal, decimal.Decimal, error) {
	data, err := ReadStorage(p, "System", "Account", hash, spec, util.TrimHex(accountID))
	if err == nil {
		var account model.AccountData
		log.Info("Account data: %v", data.ToString())
		data.ToAny(&account)
		return account.Data.Free.Add(account.Data.Reserved), decimal.Zero, nil
	}

	return decimal.Zero, decimal.Zero, err
}

func GetBalanceFromNetwork(d *dao.Dao, address string) (decimal.Decimal, error) {
	c := context.TODO()
	m, err := d.GetMetadata(c)
	if err != nil {
		log.Error("GetBalanceFromNetwork error %v", err)
		return decimal.Zero, err
	}
	spec, err := strconv.Atoi(m["specVersion"])
	if err != nil {
		log.Error("GetBalanceFromNetwork error %v", err)
		return decimal.Zero, err
	}
	balance, _, err := getFreeBalance(nil, address, "", spec)
	if err != nil {
		log.Error("GetBalanceFromNetwork error %v", err)
		return decimal.Zero, err
	}
	return balance, nil
}
