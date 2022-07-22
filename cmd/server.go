package main

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"github.com/alexpts/go-next/next"
)

func main() {
	app := next.NewApp()
	app.Use(next.Config{}, generalHandler)
	// adapter for http.HandlerFunc
	app.Use(next.Config{
		`Path`:     `/http-handle-func`,
		`Priority`: 100,
	}, next.FromHttpHandlerFunc(netHttpHandlerFunc))

	server := &fasthttp.Server{
		Handler:                       app.FasthttpHandler,
		NoDefaultDate:                 true,
		NoDefaultContentType:          true,
		NoDefaultServerHeader:         true,
		TCPKeepalive:                  true,
		GetOnly:                       true,
		DisableHeaderNamesNormalizing: true,
	}

	_ = server.ListenAndServe(":3000")
}

func generalHandler(ctx *next.HandlerCxt) {
	ctx.Response.AppendBodyString(`Ok`)
}

func netHttpHandlerFunc(w http.ResponseWriter, request *http.Request) {
	_, _ = w.Write([]byte(request.RequestURI))
}
