package next

import (
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/alexpts/go-next/next/internal"
)

const defaultRestrict = `[^/]+`

type RequestResolver struct{}

var regExpPlaceholders *regexp.Regexp

func initRegExp() *regexp.Regexp {
	if regExpPlaceholders == nil {
		regExpPlaceholders = regexp.MustCompile(`(?mU){(?P<name>.*)(:(?P<restrict>.*))?}`)
	}

	return regExpPlaceholders
}

func (r *RequestResolver) MakeRegExp(l *Layer) *regexp.Regexp {
	regexpPath := l.Path
	if regexpPath == `` {
		return nil
	}

	re := initRegExp()
	nameIndex := re.SubexpIndex(`name`)
	restrictIndex := re.SubexpIndex(`restrict`)

	matched := re.FindAllStringSubmatch(regexpPath, -1)

	for _, match := range matched {
		name := match[nameIndex]
		restrict, ok := l.Restrictions[name]

		if !ok {
			restrict = match[restrictIndex]
		}
		if restrict == `` {
			restrict = defaultRestrict
		}

		replace := "(?P<" + name + ">" + restrict + ")"
		regexpPath = strings.Replace(regexpPath, match[0], replace, -1)
	}

	return regexp.MustCompile(`(?mU)^` + regexpPath + `$`)
}

func (r *RequestResolver) ForRequest(
	l *Layer,
	request *fasthttp.RequestCtx,
	checkMethod bool,
	uriParams *UriParamsMap,
) *Layer {

	if checkMethod && !isAllowMethod(l, request) {
		return nil
	}

	if l.Path == `` {
		return l
	}

	return matchRegexpLayer(l, request, uriParams)
}

func isAllowMethod(l *Layer, req *fasthttp.RequestCtx) bool {
	if len(l.Methods) == 0 {
		return true
	}

	return internal.InSlice(l.Methods, string(req.Method()))
}

func matchRegexpLayer(l *Layer, req *fasthttp.RequestCtx, uriParams *UriParamsMap) *Layer {
	uri := req.URI().Path()
	matched := l.RegExp.FindStringSubmatch(string(uri))

	if len(matched) == 0 {
		return nil
	}

	groups := l.RegExp.SubexpNames()
	for i, name := range groups {
		if name != `` {
			(*uriParams)[name] = matched[i]
		}
	}

	return l
}
