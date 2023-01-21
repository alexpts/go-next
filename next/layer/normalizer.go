package layer

import (
	"strconv"
)

type StdNormalizer struct {
	increment   int64
	RegExpMaker IRegExpMaker
}

// Normalize -
func (n *StdNormalizer) Normalize(layer Layer) Layer {
	if layer.Name == `` {
		layer.Name = `l-` + strconv.FormatInt(n.increment, 10)
		n.increment++
	}

	layer.RegExp = n.RegExpMaker.MakeRegExp(layer)
	return layer
}
