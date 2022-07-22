package next

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func createRequest(method string, path string) *fasthttp.RequestCtx {
	r := &fasthttp.RequestCtx{
		Request: fasthttp.Request{},
	}

	r.Request.SetRequestURI(path)
	r.Request.Header.SetMethod(method)

	return r
}

func TestMinimalApp(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.Use(Config{}, func(ctx *HandlerCxt) {
		ctx.Response.AppendBodyString(`Hello`)
	})

	app.FasthttpHandler(request)
	assert.Equal(t, `Hello`, string(request.Response.Body()))
}

func TestMultiHandler(t *testing.T) {
	request := createRequest(`GET`, `/users/12/`)
	app := NewApp()

	app.Use(Config{},
		func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`Hello`)
			ctx.Next()
		},
		func(ctx *HandlerCxt) {
			_, _ = ctx.WriteString(` World`)
		},
	)
	_, err := app.Handle(request)

	assert.Nil(t, err)
	assert.Equal(t, `Hello World`, string(request.Response.Body()))
}

func TestFilterByHttpMethod(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.
		Use(Config{}, func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`1-`)
			ctx.Next()
		}).
		// Disallow POST
		Use(Config{`Methods`: `POST`}, func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`2-`)
			ctx.Next()
			ctx.Response.AppendBodyString(`2_2-`)
		}).
		// Allow one of GET
		Use(Config{`Methods`: `POST|GET`}, func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`3-`)
			ctx.Next()
		}).
		Use(Config{`Methods`: `GET`}, func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`4-`)
		})

	_, err := app.Handle(request)

	assert.Nil(t, err)
	assert.Equal(t, `1-3-4-`, string(request.Response.Body()))
}

func TestLayerPriority(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.
		Use(Config{`Priority`: 100}, func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`1-`)
		}).
		Use(Config{`Priority`: 200}, func(ctx *HandlerCxt) {
			ctx.Response.AppendBodyString(`2-`)
			ctx.Next()
		})

	_, err := app.Handle(request)

	assert.Nil(t, err)
	assert.Equal(t, `2-1-`, string(request.Response.Body()))
}

func TestDelegateToNotDefinedLayer(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.Use(Config{}, func(ctx *HandlerCxt) {
		ctx.Next()
	})

	_, err := app.Handle(request)
	assert.NotNil(t, err)
	assert.Equal(t, "can`t delegate to layer by index 1", err.Error())
}

func TestBubblePanic(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.Use(Config{}, func(ctx *HandlerCxt) {
		panic(`some panic`)
	})

	handler := func() {
		_, _ = app.Handle(request)
	}

	assert.Panics(t, handler, `some panic`)
}

func TestInterceptThrow(t *testing.T) {
	request := createRequest(`GET`, `/`)
	app := NewApp()

	app.Use(Config{}, func(ctx *HandlerCxt) {
		err := fmt.Errorf("some error")
		ctx.Panic(err, map[string]interface{}{})
	})

	_, err := app.Handle(request)
	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestFilterByPath(t *testing.T) {
	request := createRequest(`GET`, `/admin/`)
	app := NewApp()

	app.
		Use(Config{`Path`: `/users/`}, func(ctx *HandlerCxt) {
			ctx.SetUserValue(`layer1`, `users`)
			ctx.Next()
		}).
		Use(Config{`Path`: `/admin/`}, func(ctx *HandlerCxt) {
			ctx.SetUserValue(`layer2`, `admin`)
			ctx.Next()
		}).
		Use(Config{}, func(ctx *HandlerCxt) {
			ctx.SetUserValue(`layer3`, `all`)
		})

	app.FasthttpHandler(request)
	assert.Equal(t, nil, request.UserValue(`layer1`))
	assert.Equal(t, `admin`, request.UserValue(`layer2`))
	assert.Equal(t, `all`, request.UserValue(`layer3`))
}

func TestMatchUrlParam(t *testing.T) {
	request := createRequest(`GET`, `/city/london/`)
	app := NewApp()

	app.Use(Config{`Path`: `/city/{slug}/`}, func(ctx *HandlerCxt) {
		uid, ok := ctx.UriParams()["slug"]
		if ok {
			ctx.Response.AppendBodyString(uid)
		}
	})

	_, err := app.Handle(request)
	assert.Nil(t, err)
	assert.Equal(t, "london", string(request.Response.Body()))
}

func TestMatchRestricted(t *testing.T) {
	request := createRequest(`GET`, `/users/alex24/`)
	app := NewApp()

	app.Use(
		Config{
			`Name`:         `Not match regexp restriction`,
			`Path`:         `/users/{nick}/`,
			`Restrictions`: Restrictions{`nick`: `[a-z]+`},
		},
		func(ctx *HandlerCxt) {
			params := ctx.UriParams()
			uid, ok := params["nick"]
			if ok {
				ctx.Response.AppendBodyString(`1-` + uid)
			}
		})

	app.Use(
		Config{
			`Name`:         `Match regexp restriction`,
			`Path`:         `/users/{nick}/`,
			`Restrictions`: Restrictions{`nick`: `[a-z0-9]+`}, // alex24 match -> active layer
		},
		func(ctx *HandlerCxt) {
			uid, ok := ctx.UriParams()["nick"]
			if ok {
				ctx.Response.AppendBodyString(`2-` + uid)
			}
		})

	_, err := app.Handle(request)
	assert.Nil(t, err)
	assert.Equal(t, "2-alex24", string(request.Response.Body()))
}

func TestCustomStore(t *testing.T) {
	request := createRequest(`GET`, `/`)
	store := NewStore(nil, nil, nil)
	app := New(nil, store, nil)

	app.Use(Config{}, func(ctx *HandlerCxt) {
		ctx.Response.AppendBodyString(`Hello`)
	})
	_, err := app.Handle(request)

	assert.Nil(t, err)
	assert.Equal(t, `Hello`, string(request.Response.Body()))
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

			handler := func(ctx *HandlerCxt) {
				ctx.Response.AppendBodyString(string(ctx.Method()))
			}

			switch provider.method {
			case `GET`:
				app.Get(`/`, Config{}, handler)
			case `POST`:
				app.Post(`/`, Config{}, handler)
			case `PUT`:
				app.Put(`/`, Config{}, handler)
			case `PATCH`:
				app.Patch(`/`, Config{}, handler)
			case `DELETE`:
				app.Delete(`/`, Config{}, handler)
			}

			_, err := app.Handle(request)
			assert.Nil(t, err)
			assert.Equal(t, provider.expected, string(request.Response.Body()))
		})
	}
}

func TestMount(t *testing.T) {
	apiV1 := NewApp()

	apiV1.Use(Config{}, func(ctx *HandlerCxt) {
		ctx.Next()
	})
	apiV1.Get(`/users/`, Config{}, func(ctx *HandlerCxt) {
		ctx.Response.AppendBodyString(`v1 - users`)
	})

	apiV2 := NewApp()
	apiV2.Get(`/users/`, Config{}, func(ctx *HandlerCxt) {
		ctx.Response.AppendBodyString(`v2 - users`)
	})

	reuseApp := NewApp()
	reuseApp.Get(`/users/`, Config{}, func(ctx *HandlerCxt) {
		ctx.Response.AppendBodyString(`reuse - users`)
	})

	app := NewApp()
	app.
		Mount(apiV1, `/api/v1`).
		Mount(apiV2, `/api/v2`).
		Mount(reuseApp, ``)

	request := createRequest(`GET`, `/api/v1/users/`)
	app.FasthttpHandler(request)
	assert.Equal(t, `v1 - users`, string(request.Response.Body()))

	request = createRequest(`GET`, `/api/v2/users/`)
	app.FasthttpHandler(request)
	assert.Equal(t, `v2 - users`, string(request.Response.Body()))

	request = createRequest(`GET`, `/users/`)
	app.FasthttpHandler(request)
	assert.Equal(t, `reuse - users`, string(request.Response.Body()))
}
