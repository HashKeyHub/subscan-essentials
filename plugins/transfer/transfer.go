package main

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/plugins/transfer/http"
	model2 "github.com/itering/subscan/plugins/transfer/model"
	"github.com/itering/subscan/plugins/transfer/service"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/shopspring/decimal"

	"github.com/go-kratos/kratos/pkg/log"
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

	log.Info("Processing block %d, %v", block.BlockNum, extrinsic)
	var paramsInstant []storage.ExtrinsicParam
	util.UnmarshalToAnything(&paramsInstant, extrinsic.Params)
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
		if param.Type == "Address" {
			t.To = param.Value.(string)
		}
		if param.Type == "Compact<Balance>" && param.Name == "value" {
			t.Amount, _ = decimal.NewFromString(param.Value.(string))
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
