package next

import (
	"regexp"

	"github.com/valyala/fasthttp"
)

type UriParamsMap map[string]string
type Config map[string]interface{}

type Handler func(*HandlerCxt)

type FactoryContract interface {
	Create([]Handler, Config) *Layer
	CreateFromConfig(Config) (*Layer, error)
}

// NormalizeContract - интерфейс для нормализации слоя
type NormalizeContract interface {
	Normalize(*Layer, ResolverContract)
}

type ResolverContract interface {
	MakeRegExp(*Layer) *regexp.Regexp

	ForRequest(
		*Layer,
		*fasthttp.RequestCtx,
		bool,
		*UriParamsMap, // side effect by link
	) *Layer
}
