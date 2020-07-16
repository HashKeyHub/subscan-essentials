package main

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan-plugin/storage"
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

func (a *Transfer) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}
func (a *Transfer) InitHttp(e *bm.Engine) {
	a.e = e
	http.Router(srv, a.e)
}

func (a *Transfer) Http() error {
	http.Router(srv, a.e)
	return nil
}

func (a *Transfer) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, event []storage.Event) error {
	if extrinsic.CallModule != "balances" || extrinsic.CallModuleFunction != "transfer" {
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

func (a *Transfer) Migrate() {
	db := a.d.DB()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model2.Transfer{},
	)
	db.Model(model2.Transfer{}).AddUniqueIndex("extrinsic_index", "extrinsic_index")
	db.Model(model2.Transfer{}).AddIndex("idx_from", "addr_from")
	db.Model(model2.Transfer{}).AddIndex("idx_to", "addr_to")
}
