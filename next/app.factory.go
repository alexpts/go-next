package next

func NewApp() App {
	return New(nil, nil, nil)
}

func New(
	resolver ResolverContract,
	store *LayersStore,
	runner *Runner,
) App {
	app := App{&MicroApp{
		LayersStore: store,
		resolver:    resolver,
		runner:      runner,
	}}

	if app.resolver == nil {
		app.resolver = new(RequestResolver)
	}
	if app.LayersStore == nil {
		app.LayersStore = NewStore(nil, app.resolver, nil)
	}
	if app.runner == nil {
		app.runner = NewRunner(app.resolver, nil)
	}

	return app
}
