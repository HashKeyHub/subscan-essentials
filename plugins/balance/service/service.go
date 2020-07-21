package service

import (
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/balance/dao"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/shopspring/decimal"
)

type Service struct {
	d storage.Dao
}

func (s *Service) GetAccountListJson(page, row int, order, field string, queryWhere ...string) ([]*model.Account, int) {
	return dao.GetAccountList(s.d.DB(), page, row, order, field, queryWhere...)
}

func (s *Service) GetAccount(account string) (decimal.Decimal, error) {
	return dao.GetBalanceFromNetwork(account)
}

func New(d storage.Dao) *Service {
	return &Service{
		d: d,
	}
}
