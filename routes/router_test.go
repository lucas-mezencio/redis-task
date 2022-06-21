package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"redis-task/handlers"
	"testing"
)

func setupTestRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	return r
}

func TestNoRoute(t *testing.T) {
	r := setupTestRoutes()

	t.Run("empty route '/'", func(t *testing.T) {
		r.GET("/", handlers.NoRouteHandler)

		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code, "expected not found")
	})
	t.Run("random route", func(t *testing.T) {
		r.NoRoute(handlers.NoRouteHandler)

		req, _ := http.NewRequest("GET", "/asdfasdf", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code, "expected not found")
	})
	t.Run("random route and POST", func(t *testing.T) {
		r.NoRoute(handlers.NoRouteHandler)

		req, _ := http.NewRequest("POST", "/asdfasdf", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code, "expected not found")
	})
}
