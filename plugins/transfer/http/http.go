package http

import (
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	"github.com/itering/subscan/plugins/transfer/service"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
	"strings"
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
			s.POST("transfers", transfers)
		}
	}
}

func transfers(c *bm.Context) {
	p := new(struct {
		Row        int      `json:"row" validate:"min=1,max=1000"`
		Page       int      `json:"page" validate:"min=0"`
		Order      string   `json:"order" validate:"omitempty,oneof=desc asc"`
		OrderField string   `json:"order_field" validate:"omitempty"`
		Address    string   `json:"address" validate:"omitempty"`
		Types      []string `json:"types" validate:"omitempty"`
	})
	if err := c.BindWith(p, binding.JSON); err != nil {
		return
	}
	var query []string

	log.Info("api/scan/transfers params: %v", p)

	if p.Address != "" {
		a := ss58.Decode(p.Address, util.StringToInt(util.AddressType))
		var q []string
		for _, t := range p.Types {
			if t == "RECEIVE" {
				q = append(q, fmt.Sprintf("addr_to = '%s'", a))
			}
			if t == "TRANSFER" {
				q = append(q, fmt.Sprintf("addr_from = '%s'", a))
			}
		}
		query = append(query, strings.Join(q, " or "))
	}
	if p.OrderField == "" {
		p.OrderField = "block_num"
	}

	list, count := svc.GetTransferListJson(p.Page, p.Row, p.Order, p.OrderField, query...)
	c.JSON(map[string]interface{}{
		"list": list, "count": count,
	}, nil)
}
