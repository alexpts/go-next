package resolver

import "github.com/alexpts/go-next/next/layer"

type IResolver interface {
	ForRequest(
		*layer.Layer,
		*layer.HandlerCtx,
		bool,
	) *layer.Layer
}
