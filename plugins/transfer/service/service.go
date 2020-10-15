package service

import (
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/plugins/transfer/dao"
	"github.com/itering/subscan/plugins/transfer/model"
	"github.com/itering/subscan/util/address"

	"github.com/go-kratos/kratos/pkg/log"
)

type Service struct {
	d storage.Dao
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}

func (s *Service) GetTransferListJson(page, row int, order, field string, queryWhere ...string) ([]model.Transfer, int) {
	list, count := dao.FindTransfer(s.d, page, row, order, field, queryWhere...)
	for i := range list {
		list[i].To = address.SS58Address(list[i].To)
		list[i].From = address.SS58Address(list[i].From)
	}
	return list, count
}

func (s *Service) SaveTransfer(transfer *model.Transfer) error {
	log.Info("Create transfer %v", transfer)
	return dao.SaveTransfer(s.d, transfer)
}
