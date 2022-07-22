package next

import "github.com/valyala/fasthttp"

// App is wrapper for MicroApp for add users friendly methods
type App struct {
	*MicroApp
}

// FasthttpHandler - handler for fasthttp server
func (app *App) FasthttpHandler(request *fasthttp.RequestCtx) {
	_, _ = app.Handle(request)
}

// Use - attach handlers with options to micro application
func (app *App) Use(options Config, handlers ...Handler) *App {
	layer := app.factory.Create(handlers, options)
	app.LayersStore.AddLayer(layer)
	return app
}

// Mount - mount layers from another micro application to current micro application
func (app *App) Mount(app2 App, prefix string) *App {
	for _, layer := range app2.GetLayers() {
		newLayer := *layer

		newLayer.Path = prefix + layer.Path
		if layer.Path == `` {
			newLayer.Path += `/.*`
		}

		app.AddLayer(&newLayer)
	}

	return app
}
