package layer

import (
	"github.com/valyala/fasthttp"
)

type INextHandler interface {
	GetNextHandler(*HandlerCtx) Handler
}

type HandlerCtx struct {
	*fasthttp.RequestCtx
	UriParams  map[string]string
	UserParams map[string]any
	Runner     INextHandler
}

func (s *HandlerCtx) Next() {
	handler := s.Runner.GetNextHandler(s)
	handler(s)
}
