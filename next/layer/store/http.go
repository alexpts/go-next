package store

import "github.com/alexpts/go-next/next/layer"

type Handler = layer.Handler
type Layer = layer.Layer

func (s *LayersStore) Method(method string, path string, layer Layer, handlers ...Handler) *LayersStore {
	layer.Path = path
	layer.Methods = []string{method}
	return s.Use(layer, handlers...)
}

func (s *LayersStore) Get(path string, options Layer, handlers ...Handler) *LayersStore {
	return s.Method(`GET`, path, options, handlers...)
}

func (s *LayersStore) Post(path string, options Layer, handlers ...Handler) *LayersStore {
	return s.Method(`POST`, path, options, handlers...)
}

func (s *LayersStore) Put(path string, options Layer, handlers ...Handler) *LayersStore {
	return s.Method(`PUT`, path, options, handlers...)
}

func (s *LayersStore) Patch(path string, options Layer, handlers ...Handler) *LayersStore {
	return s.Method(`PATCH`, path, options, handlers...)
}

func (s *LayersStore) Delete(path string, options Layer, handlers ...Handler) *LayersStore {
	return s.Method(`DELETE`, path, options, handlers...)
}
