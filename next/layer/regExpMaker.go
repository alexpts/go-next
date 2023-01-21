package layer

import (
	"regexp"
	"strings"
)

const defaultRestrict = `[^/]+`

var regExpPlaceholders *regexp.Regexp

type StdRegExpMaker struct{}

func initRegExp() *regexp.Regexp {
	if regExpPlaceholders == nil {
		regExpPlaceholders = regexp.MustCompile(`(?mU){(?P<name>.*)(:(?P<restrict>.*))?}`)
	}

	return regExpPlaceholders
}

func (maker *StdRegExpMaker) MakeRegExp(l Layer) *regexp.Regexp {
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
