package next

import (
	"github.com/valyala/fasthttp"
)

type PanicMessage struct {
	Error error
	Props map[string]interface{}
}

type HandlerCxt struct {
	*fasthttp.RequestCtx
	runnerContext RunnerContext
}

func NewHandlerCtx(request *fasthttp.RequestCtx, runner *Runner, layers []*Layer) HandlerCxt {
	runner.SetLayers(layers)
	runnerCtx := runner.NewRunnerContext()

	return HandlerCxt{
		RequestCtx:    request,
		runnerContext: runnerCtx,
	}
}

func (ctx *HandlerCxt) Next() {
	handler := ctx.runnerContext.GetNextHandler(ctx.RequestCtx)
	handler(ctx)
}

func (ctx *HandlerCxt) UriParams() UriParamsMap {
	return ctx.runnerContext.UriParams
}

func (ctx *HandlerCxt) Panic(error error, props map[string]interface{}) {
	panic(PanicMessage{error, props})
}
