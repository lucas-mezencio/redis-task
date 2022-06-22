package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"redis-task/database"
	"redis-task/handlers"
	"redis-task/models"
	"testing"
)

var mockBlock = models.Block{
	ID:       "C3:0",
	Name:     "Bloco teste",
	ParentID: "0",
	Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
	Value:    50000000,
}

func mockBlockOnDB() {
	db := database.ConnectWithDB()
	db.Set(database.CTX, mockBlock.ID, mockBlock, 0)
}

func unmockBlock() {
	db := database.ConnectWithDB()
	db.Del(database.CTX, mockBlock.ID)
}

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

func TestGetAllBlocksRoute(t *testing.T) {
	r := setupTestRoutes()

	t.Run("get all status code", func(t *testing.T) {
		r.GET("/blocks", handlers.GetAllBlocksHandler)

		req, _ := http.NewRequest("GET", "/blocks", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestGetBlockByIdRoute(t *testing.T) {
	r := setupTestRoutes()

	t.Run("get inexistent user", func(t *testing.T) {
		r.GET("/blocks/:id", handlers.GetBlockByIdHandler)

		req, _ := http.NewRequest("GET", "/blocks/C3", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("get user by id", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		r.GET("/users/:id", handlers.GetBlockByIdHandler)

		req, _ := http.NewRequest("GET", "/users/C3", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		var gotBlock models.Block
		err := json.Unmarshal(res.Body.Bytes(), &gotBlock)
		if err != nil {
			t.Errorf("Error %g", err)
		}
		assert.Equal(t, mockBlock, gotBlock)
	})
}
