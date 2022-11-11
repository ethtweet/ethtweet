package global

import "context"

var globalContext = NewContext(context.Background())

type Context struct {
	Ctx context.Context
	Cc  context.CancelFunc
}

func NewContext(parent context.Context) *Context {
	c, cc := context.WithCancel(parent)
	return &Context{
		Ctx: c,
		Cc:  cc,
	}
}

func GetGlobalCtx() context.Context {
	return globalContext.Ctx
}

func CancelGlobalCtx() {
	cc := globalContext.Cc
	cc()
}
