package next

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Runner struct {
	resolver ResolverContract
	layers   []*Layer
}

// RunnerContext data for 1 web transaction for 1 client
type RunnerContext struct {
	runner     *Runner
	layer      *Layer
	handlerPos int
	layerPos   int

	UriParams UriParamsMap
}

func NewRunner(resolver ResolverContract, layers []*Layer) *Runner {
	return &Runner{
		resolver: resolver,
		layers:   layers,
	}
}

func (r *Runner) NewRunnerContext() RunnerContext {
	return RunnerContext{
		UriParams: UriParamsMap{},
		runner:    r,
	}
}

func (r *Runner) SetLayers(layers []*Layer) {
	r.layers = layers
}

func (rc *RunnerContext) GetNextHandler(request *fasthttp.RequestCtx) Handler {
	if rc.layer == nil {
		layer := rc.getNextLayer(request)

		handler := layer.Handlers[rc.handlerPos]
		rc.handlerPos++
		rc.layer = layer

		return handler
	}

	if rc.handlerPos == len(rc.layer.Handlers) {
		rc.layer = nil
		return rc.GetNextHandler(request)
	}

	handler := rc.layer.Handlers[rc.handlerPos]
	rc.handlerPos++
	return handler
}

func (rc *RunnerContext) getNextLayer(request *fasthttp.RequestCtx) *Layer {
	rc.handlerPos = 0

	for rc.layerPos < len(rc.runner.layers) {
		layer := rc.runner.layers[rc.layerPos]
		layer = rc.runner.resolver.ForRequest(layer, request, true, &rc.UriParams)

		rc.layerPos++
		if layer != nil {
			return layer
		}
	}

	panic(PanicMessage{
		fmt.Errorf("can`t delegate to layer by index %d", rc.layerPos),
		nil,
	})
}
