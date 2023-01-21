package resolver

import (
	"github.com/alexpts/go-next/next/internal"
	"github.com/alexpts/go-next/next/layer"
)

type RequestResolver struct{}

func (r *RequestResolver) ForRequest(
	l *layer.Layer,
	request *layer.HandlerCtx,
	checkMethod bool,
) *layer.Layer {

	if checkMethod && !isAllowMethod(l, request) {
		return nil
	}

	if l.Path == `` {
		return l
	}

	return matchRegexpLayer(l, request) // заматченные параметры вернуть
}

func isAllowMethod(l *layer.Layer, req *layer.HandlerCtx) bool {
	if len(l.Methods) == 0 {
		return true
	}

	return internal.InSlice(l.Methods, string(req.Method()))
}

func matchRegexpLayer(l *layer.Layer, req *layer.HandlerCtx) *layer.Layer {
	uri := string(req.URI().Path())
	matched := l.RegExp.FindStringSubmatch(uri)

	if len(matched) == 0 {
		return nil
	}

	groups := l.RegExp.SubexpNames()
	for i, name := range groups {
		if name != `` {
			req.UriParams[name] = matched[i]
		}
	}

	return l
}
