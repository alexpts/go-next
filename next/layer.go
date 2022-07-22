package next

import (
	"regexp"
)

type Restrictions map[string]string

type Layer struct {
	Handlers []Handler
	Name     string
	Path     string
	RegExp   *regexp.Regexp
	Priority int
	// Context  map[string]interface{}

	Methods      []string // #уточнить до ENUM
	Restrictions Restrictions
}
