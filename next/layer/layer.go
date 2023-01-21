package layer

import (
	"regexp"
)

type Restrictions map[string]string
type Handler func(ctx *HandlerCtx)

type Layer struct {
	Handlers []Handler

	Name         string
	Path         string
	RegExp       *regexp.Regexp
	Priority     int
	Methods      []string
	Restrictions Restrictions

	Meta map[string]any
}

func (l Layer) WithHandlers(handlers ...Handler) Layer {
	l.Handlers = handlers
	return l
}
