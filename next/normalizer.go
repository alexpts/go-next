package next

import (
	"strconv"
)

type StdNormalizer struct {
	increment int64
}

// Normalize @todo create bench resolver ResolverContract and resolver *ResolverContract
func (n *StdNormalizer) Normalize(layer *Layer, resolver ResolverContract) {
	if layer.Name == `` {
		layer.Name = `l-` + strconv.FormatInt(n.increment, 10)
		n.increment++
	}

	layer.RegExp = resolver.MakeRegExp(layer)
}
