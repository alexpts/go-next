package next

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreparePlaceholderRegExp(t *testing.T) {
	type placeholderRegExpProvider struct {
		uri          string
		expected     string
		restrictions map[string]string
	}

	tests := map[string]placeholderRegExpProvider{
		"prepare city placeholder": {
			uri:          `/city/{slug}/`,
			expected:     `(?mU)^/city/(?P<slug>[^/]+)/$`,
			restrictions: map[string]string{},
		},
		"replace placeholder to named group regexp": {
			uri:          `/users/{role}/{id}/`,
			expected:     `(?mU)^/users/(?P<role>[^/]+)/(?P<id>[^/]+)/$`,
			restrictions: map[string]string{},
		},
		"with restriction regexp for placeholder": {
			uri:      `/users/{id}/`,
			expected: `(?mU)^/users/(?P<id>\d+)/$`,
			restrictions: map[string]string{
				"id": `\d+`,
			},
		},
		"with restriction regexp for placeholder #2": {
			uri:      `/users/{role}/{id}/`,
			expected: `(?mU)^/users/(?P<role>[^/]+)/(?P<id>[1-9]{1}[0-9]{0,})/$`,
			restrictions: map[string]string{
				"id": `[1-9]{1}[0-9]{0,}`,
			},
		},
		"inline restrict placeholder": {
			uri:          `/users/{id:\d+}/city/{slug}/`,
			expected:     `(?mU)^/users/(?P<id>\d+)/city/(?P<slug>[^/]+)/$`,
			restrictions: map[string]string{},
		},
	}

	for name, provider := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := RequestResolver{}
			l := Layer{Path: provider.uri, Restrictions: provider.restrictions}

			regexp := resolver.MakeRegExp(&l)
			assert.Equal(t, provider.expected, regexp.String())
		})
	}
}

func TestEmptyPath(t *testing.T) {
	resolver := RequestResolver{}
	l := Layer{Path: ``}
	assert.Nil(t, resolver.MakeRegExp(&l))
}

func TestMatchRegexpLayer(t *testing.T) {
	type MatchProvider struct {
		uri          string
		path         string
		expected     UriParamsMap
		restrictions map[string]string
	}

	tests := map[string]MatchProvider{
		"prepare city placeholder": {
			uri:  `/location/{city}/`,
			path: `/location/london/`,
			expected: UriParamsMap{
				"city": "london",
			},
			restrictions: map[string]string{},
		},
		"prepare 2 placeholder": {
			uri:  `/location/{city}/`,
			path: `/location/london/`,
			expected: UriParamsMap{
				"city": "london",
			},
			restrictions: map[string]string{},
		},
	}

	for name, provider := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := RequestResolver{}
			l := &Layer{Path: provider.uri, Restrictions: provider.restrictions}
			l.RegExp = resolver.MakeRegExp(l)

			request := createRequest(`GET`, provider.path)
			uriParams := UriParamsMap{}
			result := resolver.ForRequest(l, request, true, &uriParams)

			assert.NotNil(t, result)
			assert.Len(t, uriParams, len(provider.expected))
			for name, value := range provider.expected {
				assert.Equal(t, value, uriParams[name])
			}
		})
	}
}

func BenchmarkPreparePlaceholderRegExp(b *testing.B) {
	type placeholderRegExpProvider struct {
		uri          string
		expected     string
		restrictions map[string]string
	}

	tests := map[string]placeholderRegExpProvider{
		"prepare city placeholder": {
			uri:          `/city/{slug}/`,
			expected:     `(?mU)^/city/(?P<slug>[^/]+)/$`,
			restrictions: map[string]string{},
		},
		"replace placeholder to named group regexp": {
			uri:          `/users/{role}/{id}/`,
			expected:     `(?mU)^/users/(?P<role>[^/]+)/(?P<id>[^/]+)/$`,
			restrictions: map[string]string{},
		},
		"with restriction regexp for placeholder": {
			uri:      `/users/{id}/`,
			expected: `(?mU)^/users/(?P<id>\d+)/$`,
			restrictions: map[string]string{
				"id": `\d+`,
			},
		},
		"with restriction regexp for placeholder #2": {
			uri:      `/users/{role}/{id}/`,
			expected: `(?mU)^/users/(?P<role>[^/]+)/(?P<id>[1-9]{1}[0-9]{0,})/$`,
			restrictions: map[string]string{
				"id": `[1-9]{1}[0-9]{0,}`,
			},
		},
		"inline restrict placeholder": {
			uri:          `/users/{id:\d+}/city/{slug}/`,
			expected:     `(?mU)^/users/(?P<id>\d+)/city/(?P<slug>[^/]+)/$`,
			restrictions: map[string]string{},
		},
	}

	for name, provider := range tests {
		b.Run(name, func(b *testing.B) {
			resolver := RequestResolver{}
			l := Layer{Path: provider.uri, Restrictions: provider.restrictions}

			for i := 0; i <= b.N; i++ { // b.N
				resolver.MakeRegExp(&l)
			}
		})
	}
}

//func BenchmarkPreparePlaceholderRegExp2(b *testing.B) {
//	uri := `/city/{slug}/`
//	path := `/city/london/`
//	restrictions := map[string]string{}
//
//	resolver := RequestResolver{}
//	l := &Layer{Path: uri, Restrictions: restrictions}
//	l.RegExp = resolver.MakeRegExp(l)
//
//	request := createRequest(`GET`, path)
//	uriParams := UriParamsMap{}
//
//	for i := 0; i <= b.N; i++ {
//		resolver.ForRequest(l, request, true, &uriParams)
//	}
//}
