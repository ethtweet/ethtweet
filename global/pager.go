package global

import "github.com/kataras/iris/v12"

type Pager struct {
	Offset  int
	Limit   int
	Current int
}

const DefaultPageSize = 20

func NewPager(ctx iris.Context) *Pager {
	p := &Pager{}
	p.Current = ctx.URLParamIntDefault("page", 1)
	p.Limit = ctx.URLParamIntDefault("size", DefaultPageSize)
	if p.Current <= 0 {
		p.Current = 1
	}
	if p.Limit <= 0 {
		p.Limit = DefaultPageSize
	}
	p.Offset = (p.Current - 1) * p.Limit
	return p
}
