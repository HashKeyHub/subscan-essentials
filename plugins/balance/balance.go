package balance

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/balance/http"
	"github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/shopspring/decimal"
)

var srv *service.Service

type Account struct {
	d storage.Dao
	e *bm.Engine
}

func New() *Account {
	return &Account{}
}

func (a *Account) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *Account) InitHttp() []router.Http {
	// TODO
	return []router.Http{}
}

func (a *Account) InitHttp2(e *bm.Engine) {
	http.Router(srv, e)
}

func (a *Account) ProcessExtrinsic(block *storage.Block, extrinsic *storage.Extrinsic, events []storage.Event) error {
	return nil
}

func (a *Account) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	return nil
}

func (a *Account) SubscribeExtrinsic() []string {
	return nil
}

func (a *Account) SubscribeEvent() []string {
	return []string{"system"}
}

func (a *Account) Version() string {
	return "0.1"
}

func (a *Account) Migrate() {
	a.d.AutoMigration(&model.Account{},)
	a.d.AddUniqueIndex(&model.Account{}, "address", "address")
}
