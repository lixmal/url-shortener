package backend

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lixmal/url-shortener/pkg/database"
	"github.com/lixmal/url-shortener/pkg/token"
)

const (
	tokenLen = 10
)

// url is a wrapper used for validation/binding
type url struct {
	URL string `uri:"url" binding:"required,url,min=3,max=2048"`
}

// identifier is a wrapper used for validation/binding
type identifier struct {
	Identifier string `uri:"identifier" binding:"required,printascii,min=3,max=255"`
}

func shorten(c *gin.Context) {
	var boundURL url
	if err := c.ShouldBindUri(&boundURL); err != nil {
		bindError(c, err)
		return
	}

	id := token.New(tokenLen)

	if _, ok := database.Lookup(id); ok {
		// Collision. Unlikely, but can happen with short identifiers.
		// Could regenerate the id with a few retries, but for simplicity we'll just log and report the error
		log.Printf("Collision detected for id: %s", id)

		// TODO: Use c.Negotiate to send response based on Accept header
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "unable to store shortened url, please retry",
			},
		)
	} else {
		database.Set(id, boundURL.URL)

		c.JSON(
			// StatusCreated would fit better here
			http.StatusOK,
			gin.H{
				"shortened_url": getFullURL(c.Request, id),
			},
		)
	}
}

func lookup(c *gin.Context) {
	var id identifier
	if err := c.ShouldBindUri(&id); err != nil {
		bindError(c, err)
		return
	}

	if shortenedURL, ok := database.Lookup(id.Identifier); ok {
		c.JSON(
			http.StatusOK,
			gin.H{
				"url": shortenedURL,
			},
		)
	} else {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "given identifier was not found",
			},
		)
	}
}

func redirect(c *gin.Context) {
	var id identifier
	if err := c.ShouldBindUri(&id); err != nil {
		c.Status(404)
		return
	}

	if shortenedURL, ok := database.Lookup(id.Identifier); ok {
		// Temporary because we might reuse identifiers at some point
		c.Redirect(http.StatusTemporaryRedirect, shortenedURL)
	} else {
		c.Status(404)
	}
}

func getFullURL(req *http.Request, urlPath string) string {
	/*
	* Reconstructing the full url on the server is cumbersome and should rather be done by the client.
	* Although not required, we try to guess the scheme here.
	* But it could also be ws, wss or something else.
	* In case of a reverse proxy in front we'd also need it to rewrite our URLs or provide
	* the external scheme+host somewhere in the configuration here.
	*/

	url := *req.URL
	url.Path = urlPath
	url.RawQuery = ""

	if !url.IsAbs() {
		url.Host = req.Host
		if req.TLS == nil {
			url.Scheme = "http"
		} else {
			url.Scheme = "https"
		}
	}
	return url.String()
}
