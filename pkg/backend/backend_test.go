package backend_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	stdurl "net/url"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lixmal/url-shortener/pkg/backend"
)

func init() {
	// Set up test environment
	gin.SetMode(gin.TestMode)
}

// TestBackend tests the whole API, happy path only
func TestBackend(t *testing.T) {
	router := backend.SetupRoutes()
	w := httptest.NewRecorder()

	// Request shortened url
	req, err := http.NewRequest("POST", "http://localhost/api/v1/shorten/https:%2F%2Fgoogle.com", nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Regexp(t, regexp.MustCompile(`^{"shortened_url":"http://localhost/[a-zA-Z0-9]{10}"}$`), w.Body.String())

	// Parse JSON Body
	shortenedUrl := struct {
		ShortenedUrl string `json:"shortened_url"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &shortenedUrl)
	require.NoError(t, err)

	// Parse URL
	parsedUrl, err := stdurl.Parse(shortenedUrl.ShortenedUrl)
	require.NoError(t, err)

	// Lookup shortened URL
	w.Body.Reset()
	req, err = http.NewRequest("GET", "http://localhost/api/v1/lookup"+parsedUrl.Path, nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"url":"https://google.com"}`, w.Body.String())

	// Call shortened URL for redirection
	w.Body.Reset()
	req, err = http.NewRequest("GET", parsedUrl.Path, nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "https://google.com", w.Header().Get("Location"))
}
