package dao

import (
	"fmt"
	"strings"

	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/plugins/transfer/model"
)

func SaveTransfer(db storage.Dao, m *model.Transfer) error {
	return db.Create(m)
}

func FindTransfer(db storage.Dao, page, row int, order, field string, where ...string) ([]model.Transfer, int) {
	var t []model.Transfer
	option := storage.Option {
		Page: page,
		PageSize: row,
		Order: fmt.Sprintf("%s %s", field, order),
	}

	count, ret := db.FindBy(&t, strings.Join(where, " "), &option)
	if count == 0 || ret != true {
		return t, 0
	}

	return t, count
}
