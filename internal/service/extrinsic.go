package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/internal/service/transaction"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"strings"
)

func (s *Service) createExtrinsic(c context.Context,
	txn *dao.GormDB,
	block *model.ChainBlock,
	encodeExtrinsics []string,
	decodeExtrinsics []map[string]interface{},
	eventMap map[string][]model.ChainEvent,
	finalized bool,
	spec int,
) (int, int, map[string]string, map[string]decimal.Decimal, error) {

	var (
		blockTimestamp int
		e              []model.ChainExtrinsic
		err            error
	)
	extrinsicFee := make(map[string]decimal.Decimal)

	eb, _ := json.Marshal(decodeExtrinsics)
	_ = json.Unmarshal(eb, &e)

	hash := make(map[string]string)

	s.dao.DropExtrinsicNotFinalizedData(c, block.BlockNum, finalized)

	for index, extrinsic := range e {
		extrinsic.CallModule = strings.ToLower(extrinsic.CallModule)
		extrinsic.BlockNum = block.BlockNum
		extrinsic.ExtrinsicIndex = fmt.Sprintf("%d-%d", extrinsic.BlockNum, index)
		extrinsic.Success = s.getExtrinsicSuccess(eventMap[extrinsic.ExtrinsicIndex])
		extrinsic.Finalized = finalized

		s.getTimestamp(&extrinsic)
		if extrinsic.BlockTimestamp != 0 {
			blockTimestamp = extrinsic.BlockTimestamp
		} else {
			extrinsic.BlockTimestamp = blockTimestamp
		}

		if extrinsic.ExtrinsicHash != "" {
			extrinsic.Fee = transaction.GetExtrinsicFee(encodeExtrinsics[index])
			extrinsicFee[extrinsic.ExtrinsicIndex] = extrinsic.Fee
			hash[extrinsic.ExtrinsicIndex] = extrinsic.ExtrinsicHash
		}

		if err = s.dao.CreateExtrinsic(c, txn, &extrinsic); err == nil {
			s.afterExtrinsic(block, &extrinsic, eventMap[extrinsic.ExtrinsicIndex])
		} else {
			return 0, 0, nil, nil, err
		}
	}
	return len(e), blockTimestamp, hash, extrinsicFee, err
}

func (s *Service) getTimestamp(extrinsic *model.ChainExtrinsic) {
	if extrinsic.CallModule != "timestamp" {
		return
	}

	var paramsInstant []model.ExtrinsicParam
	util.UnmarshalToAnything(&paramsInstant, extrinsic.Params)

	for _, p := range paramsInstant {
		if p.Name == "now" {
			extrinsic.BlockTimestamp = util.IntFromInterface(p.Value)
		}
	}
}

func (s *Service) getExtrinsicSuccess(e []model.ChainEvent) bool {
	f := false
	for _, event := range e {
		if strings.EqualFold(event.ModuleId, "system") {
			f = strings.EqualFold(event.EventId, "ExtrinsicFailed")
			if f {
				break
			}
		}
	}
	return !f
}

func (s *Service) afterExtrinsic(block *model.ChainBlock, extrinsic *model.ChainExtrinsic, events []model.ChainEvent) {
	block.BlockTimestamp = extrinsic.BlockTimestamp
	pBlock := block.AsPluginBlock()
	pExtrinsic := extrinsic.AsPluginExtrinsic()

	var pEvents []storage.Event
	for _, event := range events {
		pEvents = append(pEvents, *event.AsPluginEvent())
	}

	for _, plugin := range plugins.RegisteredPlugins {
		_ = plugin.ProcessExtrinsic(pBlock, pExtrinsic, pEvents)
	}
}
