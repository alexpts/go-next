package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInSlice(t *testing.T) {
	type InSliceProvider struct {
		v        string
		slice    []string
		expected bool
	}

	tests := map[string]InSliceProvider{
		`In slice`:     {`GET`, []string{`GET`, `POST`}, true},
		`In slice #2`:  {`GET`, []string{`POST`, `GET`}, true},
		`Not in slice`: {`DELETE`, []string{`GET`, `POST`}, false},
		`Empty slice`:  {`GET`, []string{}, false},
	}

	for name, provider := range tests {
		t.Run(name, func(t *testing.T) {
			actual := InSlice(provider.slice, provider.v)
			assert.Equal(t, provider.expected, actual)
		})
	}
}
