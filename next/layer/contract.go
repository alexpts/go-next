package layer

import "regexp"

type IRegExpMaker interface {
	MakeRegExp(Layer) *regexp.Regexp
}

type INormalizer interface {
	Normalize(layer Layer) Layer
}
