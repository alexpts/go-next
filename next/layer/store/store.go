package store

import (
	"sort"

	"github.com/alexpts/go-next/next/layer"
)

type LayersStore struct {
	normalizer layer.INormalizer
	layers     []layer.Layer
	sorted     bool
}

func New(normalizer layer.INormalizer) *LayersStore {
	return &LayersStore{
		normalizer: normalizer,
	}
}

func (s *LayersStore) AddLayer(layer layer.Layer) *LayersStore {
	nLayer := s.normalizer.Normalize(layer)
	s.layers = append(s.layers, nLayer)
	s.sorted = false
	return s
}

func (s *LayersStore) Use(layer layer.Layer, handlers ...layer.Handler) *LayersStore {
	return s.AddLayer(
		layer.WithHandlers(handlers...),
	)
}

func (s *LayersStore) GetLayers() []layer.Layer {
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
