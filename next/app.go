package next

import (
	"github.com/valyala/fasthttp"

	"github.com/alexpts/go-next/next/layer"
	"github.com/alexpts/go-next/next/layer/resolver"
	"github.com/alexpts/go-next/next/layer/store"
	"github.com/alexpts/go-next/next/runner"
)

type MicroApp struct {
	*store.LayersStore
	resolver      resolver.IResolver
	fallbackLayer *layer.Layer
}

func NewApp() MicroApp {
	return ProvideMicroApp(nil, nil, nil)
}

func ProvideMicroApp(
	resolverObj resolver.IResolver,
	storeObj *store.LayersStore,
	fallbackLayer *layer.Layer,
) MicroApp {
	if storeObj == nil {
		storeObj = store.New(
			&layer.StdNormalizer{
				RegExpMaker: &layer.StdRegExpMaker{},
			},
		)
	}

	if resolverObj == nil {
		resolverObj = &resolver.RequestResolver{}
	}

	app := MicroApp{
		LayersStore:   storeObj,
		resolver:      resolverObj,
		fallbackLayer: fallbackLayer,
	}

	return app
}

// Mount - mount layers from another micro application to current micro application
func (app *MicroApp) Mount(app2 MicroApp, prefix string) *MicroApp {
	for _, l := range app2.GetLayers() {
		newLayer := l

		newLayer.Path = prefix + l.Path
		if l.Path == `` {
			newLayer.Path += `/.*`
		}

		app.AddLayer(newLayer)
	}

	return app
}

func (app *MicroApp) Handle(req *fasthttp.RequestCtx) *layer.HandlerCtx {
	nextCtx := &layer.HandlerCtx{
		RequestCtx: req,
		Runner:     app.createRunner(),
		UriParams:  make(map[string]string),
		UserParams: make(map[string]any),
	}

	nextCtx.Next()
	return nextCtx
}

func (app *MicroApp) FastHttpHandler(req *fasthttp.RequestCtx) {
	_ = app.Handle(req)
}

func (app *MicroApp) createRunner() layer.INextHandler {

	return &runner.Runner{
		Layers:        app.GetLayers(),
		Resolver:      &resolver.RequestResolver{},
		FallbackLayer: app.fallbackLayer,
	}
}
