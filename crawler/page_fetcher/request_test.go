package page_fetcher

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidity(t *testing.T) {
	{
		req := Request{}
		resp := http.Response{}
		assert.True(t, req.acceptableResponse(&resp))
	}
	{
		req := Request{
			URL:          nil,
			HTTPReferrer: "",

			AcceptableContentTypes: map[string]struct{}{
				"text/html": {},
			},
		}
		resp := http.Response{}
		assert.False(t, req.acceptableResponse(&resp))
	}
	{
		req := Request{
			URL:          nil,
			HTTPReferrer: "",
			AcceptableContentTypes: map[string]struct{}{
				"text/html": {},
			},
		}
		headers := http.Header{}
		headers.Add("Content-Type", "text/css")
		resp := http.Response{
			Header: headers,
		}
		assert.False(t, req.acceptableResponse(&resp))
	}
	{
		req := Request{
			URL:          nil,
			HTTPReferrer: "",
			AcceptableContentTypes: map[string]struct{}{
				"text/html": {},
			},
		}
		headers := http.Header{}
		headers.Add("Content-Type", "text/html")
		resp := http.Response{
			Header: headers,
		}
		assert.True(t, req.acceptableResponse(&resp))
	}
	{
		req := Request{
			URL:          nil,
			HTTPReferrer: "",
			AcceptableContentTypes: map[string]struct{}{
				"text/html": {},
			},
		}
		headers := http.Header{}
		headers.Add("Content-Type", "text/html; charset=utf-8")
		resp := http.Response{
			Header: headers,
		}
		assert.True(t, req.acceptableResponse(&resp))
	}
}
