package next

func (s *LayersStore) Method(method string, path string, options Config, handlers ...Handler) *LayersStore {
	options[`Path`] = path
	options[`Methods`] = method

	layer := s.factory.Create(handlers, options)
	return s.AddLayer(layer)
}

func (s *LayersStore) Get(path string, options Config, handlers ...Handler) *LayersStore {
	return s.Method(`GET`, path, options, handlers...)
}

func (s *LayersStore) Post(path string, options Config, handlers ...Handler) *LayersStore {
	return s.Method(`POST`, path, options, handlers...)
}

func (s *LayersStore) Put(path string, options Config, handlers ...Handler) *LayersStore {
	return s.Method(`PUT`, path, options, handlers...)
}

func (s *LayersStore) Patch(path string, options Config, handlers ...Handler) *LayersStore {
	return s.Method(`PATCH`, path, options, handlers...)
}

func (s *LayersStore) Delete(path string, options Config, handlers ...Handler) *LayersStore {
	return s.Method(`DELETE`, path, options, handlers...)
}
