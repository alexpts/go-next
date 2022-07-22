package next

import (
	"fmt"
	"strings"
)

type Factory struct{}

// Create - create new Layer
func (f *Factory) Create(handlers []Handler, config Config) *Layer {
	layer := Layer{Handlers: handlers}

	for name, value := range config {
		switch name {
		case `Name`:
			layer.Name = value.(string)
		case `Priority`:
			layer.Priority = value.(int)
		case `Path`:
			layer.Path = value.(string)
		case `Methods`:
			methods := value.(string)
			if methods != `` {
				layer.Methods = strings.Split(methods, `|`) // GET|POST|PUT
			}
		case `Restrictions`:
			layer.Restrictions = value.(Restrictions)
		}
	}

	return &layer
}

// CreateFromConfig - crate layer from declarative config
func (f *Factory) CreateFromConfig(config Config) (*Layer, error) {
	value, ok := config["Handlers"]
	if ok == false {
		return nil, fmt.Errorf("invalid config")
	}

	handlers := value.([]Handler)
	return f.Create(handlers, config), nil
}
