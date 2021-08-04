package backend

import (
	"net/http"
	stdurl "net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullURL(t *testing.T) {
	type args struct {
		req     *http.Request
		urlPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"relative",
			args{
				&http.Request{
					Host: "somehost.tld",
					URL: &stdurl.URL{
						Path:     "bogus",
						RawQuery: "bogus",
					},
				},
				"reltest",
			},
			"http://somehost.tld/reltest",
		},

		{
			"absolute",
			args{
				&http.Request{
					Host: "somehost.tld",
					URL: &stdurl.URL{
						Scheme:   "https",
						Host:     "someotherhost.tld",
						Path:     "bogus",
						RawQuery: "bogus",
					},
				},
				"abstest",
			},
			"https://someotherhost.tld/abstest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, getFullURL(tt.args.req, tt.args.urlPath), tt.want)
		})
	}
}
