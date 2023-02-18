package next

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/alexpts/go-next/next/layer"
)

type HandlerCtx = layer.HandlerCtx
type Layer = layer.Layer

func createRequest(method string, path string) *fasthttp.RequestCtx {
	r := fasthttp.RequestCtx{
		Request: fasthttp.Request{},
	}

	r.Request.SetRequestURI(path)
	r.Request.Header.SetMethod(method)

	return &r
}

func TestMinimalApp(t *testing.T) {
	app := NewApp()
	request := createRequest(`GET`, `/`)

	app.Use(Layer{},
		func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`Hello`)
			return nil
		},
	)

	app.Handle(request)
	assert.Equal(t, `Hello`, string(request.Response.Body()))
}

func TestMultiHandler(t *testing.T) {
	app := NewApp()
	request := createRequest(`GET`, `/users/12/`)

	app.Use(Layer{},
		func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`Hello`)
			return ctx.Next()
		},
		func(ctx *HandlerCtx) error {
			_, err := ctx.WriteString(` World`)
			return err
		},
	)

	app.Handle(request)
	assert.Equal(t, `Hello World`, string(request.Response.Body()))
}

func TestMultiLayers(t *testing.T) {
	app := NewApp()
	request := createRequest(`GET`, `/users/12/`)

	app.AddLayer(Layer{}.
		WithHandlers(
			func(ctx *HandlerCtx) error {
				ctx.Response.AppendBodyString(`Hello`)
				return ctx.Next()
			},
		),
	)

	app.AddLayer(Layer{}.
		WithHandlers(
			func(ctx *HandlerCtx) error {
				_, err := ctx.WriteString(` World`)
				return err
			},
		),
	)

	app.Handle(request)
	assert.Equal(t, `Hello World`, string(request.Response.Body()))
}

func TestLayerPriority(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.
		Use(Layer{Priority: 100}, func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`1-`) // run second
			return nil
		}).
		Use(Layer{Priority: 200}, func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`2-`) // run first
			return ctx.Next()
		})

	app.Handle(request)
	assert.Equal(t, `2-1-`, string(request.Response.Body()))
}

func TestDelegateToNotDefinedLayer(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.Use(Layer{}, func(ctx *HandlerCtx) error {
		return ctx.Next()
	})

	assert.Panics(t, func() {
		app.Handle(request)
	}, "can`t delegate to layer by index 1")
}

func TestFallbackLayer(t *testing.T) {
	request := createRequest(`GET`, `/`)
	fallbackLayer := Layer{}.WithHandlers(func(ctx *HandlerCtx) error {
		ctx.Response.SetStatusCode(500)
		ctx.SetContentType("application/json")
		ctx.Response.AppendBody([]byte(`{"error": "not found handler"}`))
		return nil
	})

	app := ProvideMicroApp(nil, nil, &fallbackLayer)

	app.Use(Layer{}, func(ctx *HandlerCtx) error {
		return ctx.Next()
	})

	app.Handle(request)
	assert.Equal(t, 500, request.Response.StatusCode())
	assert.Equal(t, "application/json", string(request.Response.Header.ContentType()))
	assert.Equal(t, `{"error": "not found handler"}`, string(request.Response.Body()))
}

func TestFilterByHttpMethod(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.
		Use(Layer{}, func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`1-`)
			return ctx.Next()
		}).
		// Disallow POST
		Use(Layer{Methods: []string{`POST`}}, func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`2-`)
			_ = ctx.Next()
			ctx.Response.AppendBodyString(`2_2-`)
			return nil
		}).
		// Allow one of GET
		Use(Layer{Methods: []string{`POST`, `GET`}}, func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`3-`)
			return ctx.Next()
		}).
		Use(Layer{Methods: []string{`GET`}}, func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`4-`)
			return nil
		})

	app.Handle(request)
	assert.Equal(t, `1-3-4-`, string(request.Response.Body()))
}

func TestFilterByPath(t *testing.T) {
	request := createRequest(`GET`, `/admin/`)
	app := NewApp()

	app.
		Use(Layer{Path: `/users/`}, func(ctx *HandlerCtx) error {
			ctx.SetUserValue(`layer1`, `users`)
			return ctx.Next()
		}).
		Use(Layer{Path: `/admin/`}, func(ctx *HandlerCtx) error {
			ctx.SetUserValue(`layer2`, `admin`)
			return ctx.Next()
		}).
		Use(Layer{}, func(ctx *HandlerCtx) error {
			ctx.SetUserValue(`layer3`, `all`)
			return nil
		})

	app.Handle(request)
	assert.Equal(t, nil, request.UserValue(`layer1`))
	assert.Equal(t, `admin`, request.UserValue(`layer2`))
	assert.Equal(t, `all`, request.UserValue(`layer3`))
}

func TestMatchUrlParam(t *testing.T) {
	request := createRequest(`GET`, `/city/london/`)
	app := NewApp()

	app.Use(Layer{Path: `/city/{slug}/`}, func(ctx *HandlerCtx) error {
		uid, ok := ctx.UriParams["slug"]
		if ok {
			ctx.Response.AppendBodyString(uid)
		}

		return nil
	})

	app.Handle(request)
	assert.Equal(t, "london", string(request.Response.Body()))
}

func TestFastHttpMethod(t *testing.T) {
	type testProvider struct {
		method   string
		expected string
	}

	tests := map[string]testProvider{
		"GET": {
			method:   `GET`,
			expected: `GET`,
		},
		"POST": {
			method:   `POST`,
			expected: `POST`,
		},
		"PUT": {
			method:   `PUT`,
			expected: `PUT`,
		},
		"PATCH": {
			method:   `PATCH`,
			expected: `PATCH`,
		},
		"DELETE": {
			method:   `DELETE`,
			expected: `DELETE`,
		},
	}

	for name, provider := range tests {
		t.Run(name, func(t *testing.T) {
			request := createRequest(provider.method, `/`)
			app := NewApp()

			handler := func(ctx *HandlerCtx) error {
				ctx.Response.AppendBodyString(string(ctx.Method()))
				return nil
			}

			switch provider.method {
			case `GET`:
				app.Get(`/`, Layer{}, handler)
			case `POST`:
				app.Post(`/`, Layer{}, handler)
			case `PUT`:
				app.Put(`/`, Layer{}, handler)
			case `PATCH`:
				app.Patch(`/`, Layer{}, handler)
			case `DELETE`:
				app.Delete(`/`, Layer{}, handler)
			}

			app.Handle(request)
			assert.Equal(t, provider.expected, string(request.Response.Body()))
		})
	}
}

func TestMount(t *testing.T) {
	apiV1 := NewApp()

	apiV1.Use(Layer{}, func(ctx *HandlerCtx) error {
		return ctx.Next()
	})
	apiV1.Get(`/users/`, Layer{}, func(ctx *HandlerCtx) error {
		ctx.Response.AppendBodyString(`v1 - users`)
		return nil
	})

	apiV2 := NewApp()
	apiV2.Get(`/users/`, Layer{}, func(ctx *HandlerCtx) error {
		ctx.Response.AppendBodyString(`v2 - users`)
		return nil
	})

	reuseApp := NewApp()
	reuseApp.Get(`/users/`, Layer{}, func(ctx *HandlerCtx) error {
		ctx.Response.AppendBodyString(`reuse - users`)
		return nil
	})

	app := NewApp()
	app.
		Mount(apiV1, `/api/v1`).
		Mount(apiV2, `/api/v2`).
		Mount(reuseApp, ``)

	request := createRequest(`GET`, `/api/v1/users/`)
	app.Handle(request)
	assert.Equal(t, `v1 - users`, string(request.Response.Body()))

	request = createRequest(`GET`, `/api/v2/users/`)
	app.Handle(request)
	assert.Equal(t, `v2 - users`, string(request.Response.Body()))

	request = createRequest(`GET`, `/users/`)
	app.Handle(request)
	assert.Equal(t, `reuse - users`, string(request.Response.Body()))
}

func TestFastHttpHandler(t *testing.T) {
	app := NewApp()
	request := createRequest(`GET`, `/`)

	app.Use(Layer{},
		func(ctx *HandlerCtx) error {
			ctx.Response.AppendBodyString(`Hello`)
			return nil
		},
	)

	app.FastHttpHandler(request)
	assert.Equal(t, `Hello`, string(request.Response.Body()))
}
