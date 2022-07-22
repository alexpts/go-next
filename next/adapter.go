package next

import (
	"net/http"

	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// FromHttpHandlerFunc wraps net/http handler func to next.HandlerCxt via fasthttp request handler
func FromHttpHandlerFunc(handler http.HandlerFunc) Handler {
	fasthttpHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(handler)
	return func(cxt *HandlerCxt) {
		fasthttpHandler(cxt.RequestCtx)
	}
}
