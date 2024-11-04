package module

import (
	"net/http"
	"strings"
	"testing"

	"github.com/lingwei0604/kitty/rule/dto"
	"github.com/stretchr/testify/assert"
)

func TestDecodePayload(t *testing.T) {
	cases := []struct {
		name    string
		request *http.Request
		asserts func(t *testing.T, payload *dto.Payload)
	}{
		{
			"decode query",
			func() *http.Request {
				r, _ := http.NewRequest("GET", "http://example.org?foo=bar&foo=baz", nil)
				return r
			}(),
			func(t *testing.T, payload *dto.Payload) {
				assert.Contains(t, payload.Q["foo"], "bar")
				assert.Contains(t, payload.Q["foo"], "baz")
			},
		},
		{
			"decode body",
			func() *http.Request {
				r, _ := http.NewRequest("POST", "http://example.org", strings.NewReader(`{"foo":"bar"}`))
				return r
			}(),
			func(t *testing.T, payload *dto.Payload) {
				assert.Equal(t, payload.B["foo"], "bar")
			},
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.name, func(t *testing.T) {
			p := dto.Payload{}
			err := dto.NewDecoder().Decode(&p, cc.request)
			assert.NoError(t, err)
			cc.asserts(t, &p)
		})
	}
}
