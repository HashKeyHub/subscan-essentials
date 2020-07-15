package dao

import (
	"fmt"
	"github.com/itering/subscan/plugins/transfer/model"
	"github.com/jinzhu/gorm"
)

func SaveTransfer(db *gorm.DB, m *model.Transfer) error {
	query := db.Save(m)
	return query.Error
}

func FindTransfer(db *gorm.DB, page, row int, order, field string, where ...string) ([]model.Transfer, int) {
	var t []model.Transfer
	queryOrigin := db.Model(model.Transfer{})
	for _, w := range where {
		queryOrigin = queryOrigin.Where(w)
	}

	query := queryOrigin.Order(fmt.Sprintf("%s %s", field, order)).Offset(page * row).Limit(row).Scan(&t)
	if query == nil || query.Error != nil || query.RecordNotFound() {
		return t, 0
	}

	var count int
	queryOrigin.Count(&count)
	return t, count
}
