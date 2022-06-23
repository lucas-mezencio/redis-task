package routes

import (
	"bytes"
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
	unmockBlock()
	db := database.ConnectWithDB()
	db.Set(database.CTX, mockBlock.ID, mockBlock, 0)
}

func unmockBlock() {
	db := database.ConnectWithDB()
	db.FlushAll(database.CTX)
}

var (
	c0 = models.Block{
		ID:       "C0:0",
		Name:     "Cliente A",
		ParentID: "0",
		Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
		Value:    10000,
	}
	f1 = models.Block{
		ID:       "F1:C0",
		Name:     "FAZENDA 1",
		ParentID: "C0",
		Centroid: *geojson.NewPointGeometry([]float64{-52.9046630859375, -18.132801356084773}),
		Value:    1000,
	}
)

var treeMock = models.Tree{
	Block: c0,
	Children: []models.Tree{
		{
			Block:    f1,
			Children: nil,
		},
	},
}

func mockTree(t *testing.T) {
	unmockTree(t)
	db := database.ConnectWithDB()
	blocks := []models.Block{c0, f1}
	for _, block := range blocks {
		err := db.Set(database.CTX, block.ID, block, 0).Err()
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func unmockTree(t *testing.T) {
	db := database.ConnectWithDB()
	err := db.FlushAll(database.CTX).Err()
	if err != nil {
		t.Error(err)
	}
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

func TestGetTreeBellowIdRoute(t *testing.T) {
	r := setupTestRoutes()
	r.GET("/tree/:id", handlers.GetTreeBellowId)

	t.Run("nonexistent tree", func(t *testing.T) {
		// limpa o banco
		unmockTree(t)
		req, _ := http.NewRequest("GET", "/tree/C0", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("mocked tree", func(t *testing.T) {
		unmockTree(t)
		mockTree(t)
		req, _ := http.NewRequest("GET", "/tree/C0", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var gotTree models.Tree
		err := json.Unmarshal(res.Body.Bytes(), &gotTree)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, treeMock, gotTree)
	})
}

func TestCreateBlockRoute(t *testing.T) {
	r := setupTestRoutes()
	r.POST("/blocks", handlers.CreateBlockHandler)
	bytesUser, _ := json.Marshal(mockBlock)
	t.Run("create new block", func(t *testing.T) {
		unmockBlock()

		req, _ := http.NewRequest("POST", "/blocks", bytes.NewBuffer(bytesUser))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		var gotBlock models.Block
		err := json.Unmarshal(res.Body.Bytes(), &gotBlock)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, mockBlock, gotBlock)

		unmockBlock()
	})

	t.Run("create existing block", func(t *testing.T) {
		mockBlockOnDB()

		req, _ := http.NewRequest("POST", "/blocks", bytes.NewBuffer(bytesUser))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		unmockBlock()
	})
}

func TestUpdateBlockByIdRoute(t *testing.T) {
	r := setupTestRoutes()
	r.PUT("/blocks/:id", handlers.UpdateBlockByIdHandler)
	updatedBlock := mockBlock
	t.Run("update valid block", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		bytesUpdatedBlock, _ := json.Marshal(updatedBlock)

		req, _ := http.NewRequest("PUT", "/blocks/C3", bytes.NewBuffer(bytesUpdatedBlock))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		var gotBlock models.Block
		err := json.Unmarshal(res.Body.Bytes(), &gotBlock)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, updatedBlock, gotBlock)
	})
	t.Run("update inexistent block", func(t *testing.T) {
		unmockBlock()
		bytesUpdatedBlock, _ := json.Marshal(updatedBlock)
		req, _ := http.NewRequest("PUT", "/blocks/C3", bytes.NewBuffer(bytesUpdatedBlock))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("update invalid parent id", func(t *testing.T) {
		mockBlockOnDB()
		defer unmockBlock()
		updatedBlock.ParentID = "asdf"
		bytesUpdatedBlock, _ := json.Marshal(updatedBlock)
		req, _ := http.NewRequest("PUT", "/blocks/C3", bytes.NewBuffer(bytesUpdatedBlock))
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}

func TestDeleteByIdRoute(t *testing.T) {
	r := setupTestRoutes()
	r.DELETE("/blocks/:id", handlers.DeleteBlockByIdHandler)
	t.Run("existent block", func(t *testing.T) {
		mockBlockOnDB()

		req, _ := http.NewRequest("DELETE", "/blocks/C3", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		gotBlock := models.GetBlockById("C3")
		assert.Equal(t, models.Block{}, gotBlock)
	})

	t.Run("nonexistent block", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/blocks/C3", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("child parent key transition", func(t *testing.T) {
		mockBlockOnDB()
		child := models.Block{
			ID:       "T1:C3",
			Name:     "Bloco teste",
			ParentID: "C3",
			Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
			Value:    50000000,
		}
		err := models.CreateBlock(child)
		if err != nil {
			t.Error(err)
		}

		req, _ := http.NewRequest("DELETE", "/blocks/C3", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		updatedChild := models.GetBlockById("T1")
		assert.Equal(t, "T1:0", updatedChild.ID)
		assert.Equal(t, "0", updatedChild.ParentID)
		unmockBlock()
	})
}
