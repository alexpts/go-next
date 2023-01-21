package runner

import (
	"fmt"

	"github.com/alexpts/go-next/next/layer"
	"github.com/alexpts/go-next/next/layer/resolver"
)

type Runner struct {
	Resolver      resolver.IResolver
	Layers        []layer.Layer
	FallbackLayer *layer.Layer

	curLayer   *layer.Layer
	layerPos   int
	handlerPos int
}

func (r *Runner) GetNextHandler(request *layer.HandlerCtx) layer.Handler {
	if r.curLayer == nil {
		r.curLayer = r.getNextLayer(request)
	}

	// Если обработчики в слое закончились, то идем в следующий слой
	if r.handlerPos == len(r.curLayer.Handlers) {
		r.curLayer = nil
		return r.GetNextHandler(request)
	}

	handler := r.curLayer.Handlers[r.handlerPos]
	r.handlerPos++
	return handler
}

func (r *Runner) getNextLayer(request *layer.HandlerCtx) *layer.Layer {
	r.handlerPos = 0

	for r.layerPos < len(r.Layers) {
		refLayer := &r.Layers[r.layerPos]
		r.layerPos++

		refLayer = r.Resolver.ForRequest(refLayer, request, true)
		if refLayer != nil {
			return refLayer
		}
	}

	if r.FallbackLayer == nil {
		panic(fmt.Errorf("can`t delegate to layer by index %d", r.layerPos))
	}

	return r.FallbackLayer
}

//func (r *Runner) Reset() {
//	r.layerPos = 0
//	r.handlerPos = 0
//	r.curLayer = nil
//}
