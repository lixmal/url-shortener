// package backend implements http routing and logic for the url shortener service
package backend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// New creates a new url shortener server and returns the configured http.Server
func New(listenAddr string) *http.Server {
	return &http.Server{
		Addr:    listenAddr,
		Handler: SetupRoutes(),
	}
}

// SetupRoutes sets http routes up and returns the gin engine
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// We need to process encoded slashes
	router.UseRawPath = true

	v1 := router.Group("/api/v1")
	v1.POST("/shorten/:url", shorten)
	v1.GET("/lookup/:identifier", lookup)

	router.GET("/:identifier", redirect)

	return router
}

func bindError(c *gin.Context, err error) {
	// TODO: get FieldError and return relevant part only
	c.JSON(
		http.StatusBadRequest,
		gin.H{
			"error": err.Error(),
		},
	)
}
