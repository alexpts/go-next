package next

import (
	"sort"
)

type LayersStore struct {
	normalizer NormalizeContract
	resolver   ResolverContract
	factory    FactoryContract
	layers     []*Layer
	sorted     bool
}

func NewStore(n NormalizeContract, r ResolverContract, f FactoryContract) *LayersStore {
	store := LayersStore{
		normalizer: n,
		resolver:   r,
		factory:    f,
	}

	if store.normalizer == nil {
		store.normalizer = new(StdNormalizer)
	}
	if store.resolver == nil {
		store.resolver = new(RequestResolver)
	}
	if store.factory == nil {
		store.factory = &Factory{}
	}

	return &store
}

func (s *LayersStore) AddLayer(layer *Layer) *LayersStore {
	s.normalizer.Normalize(layer, s.resolver)
	s.layers = append(s.layers, layer)
	s.sorted = false
	return s
}

func (s *LayersStore) GetLayers() []*Layer {
	if s.sorted == false {
		s.sortByPriority()
	}

	return s.layers
}

func (s *LayersStore) sortByPriority() {
	sort.SliceStable(s.layers, func(i, j int) bool {
		return s.layers[i].Priority > s.layers[j].Priority
	})

	s.sorted = true
}
