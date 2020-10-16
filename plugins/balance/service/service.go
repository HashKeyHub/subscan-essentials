package service

import (
	id "github.com/itering/subscan/internal/dao"
	"github.com/itering/subscan/plugins/balance/dao"
	"github.com/shopspring/decimal"
)

type Service struct {
	d *id.Dao
}

func (s *Service) GetAccount(account string) (decimal.Decimal, error) {
	return dao.GetBalanceFromNetwork(s.d, account)
}

func New(d *id.Dao) *Service {
	return &Service{
		d: d,
	}
}
