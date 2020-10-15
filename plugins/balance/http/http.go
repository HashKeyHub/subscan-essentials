package http

import (
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
)

var (
	svc *service.Service
)

func Router(s *service.Service, e *bm.Engine) {
	svc = s
	g := e.Group("/api")
	{
		s := g.Group("/scan")
		{
			s.GET("accounts/:address/balance", account)
		}
	}
}

func account(c *bm.Context) {
	addr,_ := c.Params.Get("address")
	balance, _ := svc.GetAccount(ss58.Decode(addr, util.StringToInt(util.AddressType)))
	c.JSON(map[string]interface{}{
		"balance": balance,
	}, nil)
}
