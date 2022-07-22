package next

import (
	"github.com/valyala/fasthttp"
)

type MicroApp struct {
	*LayersStore
	resolver ResolverContract
	runner   *Runner
}

func catchError(recovery any) error {
	switch recovery := recovery.(type) {
	case nil:
		return nil
	case PanicMessage:
		return recovery.Error
	default:
		panic(recovery)
	}
}

func (app *MicroApp) Handle(request *fasthttp.RequestCtx) (flow HandlerCxt, err error) {
	defer func() {
		err = catchError(recover())
	}()

	flow = NewHandlerCtx(request, app.runner, app.LayersStore.GetLayers())
	flow.Next()

	return flow, nil
}
