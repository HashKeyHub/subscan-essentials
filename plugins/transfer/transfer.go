package main

import (
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/plugins/transfer/http"
	model2 "github.com/itering/subscan/plugins/transfer/model"
	"github.com/itering/subscan/plugins/transfer/service"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type Transfer struct {
	d storage.Dao
	e *bm.Engine
}

func New() *Transfer {
	return &Transfer{}
}

func (a *Transfer) InitDao(dao *dao.Dao, d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *Transfer) InitHttp() []router.Http {
	return []router.Http{}
}

func (a *Transfer) InitHttp2(e *bm.Engine) {
	a.e = e
	http.Router(srv, a.e)
}

func (a *Transfer) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, event []storage.Event) error {
	if extrinsic.CallModule != "balances" || (extrinsic.CallModuleFunction != "transfer" && extrinsic.CallModuleFunction != "transfer_keep_alive") {
		// ignore others
		return nil
	}
	log.Info("Processing block %d, %v, %s, %s", block.BlockNum, extrinsic)
	paramsInstant := extrinsic.Params.([]interface{})
	var t = model2.Transfer{
		From:           extrinsic.AccountId,
		BlockNum:       block.BlockNum,
		BlockTimestamp: block.BlockTimestamp,
		ExtrinsicIndex: extrinsic.ExtrinsicIndex,
		Fee:            extrinsic.Fee,
		Finalized:      block.Finalized,
		Success:        extrinsic.Success,
	}
	log.Info("extrinsic params: %v", paramsInstant)
	for _, param := range paramsInstant {
		p := param.(map[string]interface{})
		if p["Type"] == "Address" {
			to, err := decodeMultiAddress(p["Value"].(map[string]interface{}))
			if err != nil {
				log.Error("%s", err)
				return err
			}
			t.To = to
		}
		if p["Type"] == "Compact<Balance>" && p["Name"] == "value" {
			t.Amount, _ = decimal.NewFromString(p["Value"].(string))
		}
	}
	return srv.SaveTransfer(&t)
}

func (a *Transfer) ProcessEvent(*storage.Block, *storage.Event, decimal.Decimal) error {
	return nil
}

func (a *Transfer) SubscribeExtrinsic() []string {
	return nil
}

func (a *Transfer) SubscribeEvent() []string {
	return []string{"system"}
}

func (a *Transfer) Version() string {
	return "0.1"
}

func (a *Transfer) Migrate() {
	a.d.AutoMigration(
		&model2.Transfer{},
	)
	a.d.AddUniqueIndex(&model2.Transfer{}, "extrinsic_index", "extrinsic_index")
	a.d.AddIndex(&model2.Transfer{}, "idx_from", "addr_from")
	a.d.AddIndex(&model2.Transfer{}, "idx_to", "addr_to")
}

func decodeMultiAddress(m map[string]interface{}) (string, error) {
	id, ok := m["Id"]
	if ok {
		return id.(string), nil
	}

	return "", errors.New("invalid address")
}