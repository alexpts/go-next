package main

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/alexpts/go-next/next"
	"github.com/alexpts/go-next/next/layer"
)

func main() {
	app := next.NewApp()
	app.Use(layer.Layer{}, generalHandler)
	// adapter for http.HandlerFunc
	app.Use(layer.Layer{
		Path:     `/http-handle-func`,
		Priority: 100,
	}, fromMux(netHttpHandlerFunc))

	server := &fasthttp.Server{
		Handler:                       app.FastHttpHandler,
		NoDefaultDate:                 true,
		NoDefaultContentType:          true,
		NoDefaultServerHeader:         true,
		TCPKeepalive:                  true,
		GetOnly:                       true,
		DisableHeaderNamesNormalizing: true,
	}

	_ = server.ListenAndServe(":3000")
}

func generalHandler(ctx *layer.HandlerCtx) {
	ctx.Response.AppendBodyString(`Ok`)
}

func netHttpHandlerFunc(w http.ResponseWriter, request *http.Request) {
	_, _ = w.Write([]byte(request.RequestURI))
}

func fromMux(handler http.HandlerFunc) layer.Handler {
	fasthttpHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(handler)
	return func(cxt *layer.HandlerCtx) {
		fasthttpHandler(cxt.RequestCtx)
	}
}
