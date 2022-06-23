package scripts

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"redis-task/database"
	"redis-task/models"
)

func FlushDatabase() {
	db := database.ConnectWithDB()
	//defer db.Close()

	err := db.FlushAll(database.CTX).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func PopulateDatabase(blocks []models.Block) {
	db := database.ConnectWithDB()
	//defer db.Close()
	if blocks == nil {
		c0 := models.Block{
			ID:       "C0:0",
			Name:     "Cliente A",
			ParentID: "0",
			Centroid: *geojson.NewPointGeometry([]float64{-48.289546966552734, -18.931050694554795}),
			Value:    10000,
		}
		f1 := models.Block{
			ID:       "F1:C0",
			Name:     "FAZENDA 1",
			ParentID: "C0",
			Centroid: *geojson.NewPointGeometry([]float64{-52.9046630859375, -18.132801356084773}),
			Value:    1000,
		}
		f2 := models.Block{
			ID:       "F2:C0",
			Name:     "FAZENDA 2",
			ParentID: "C0",
			Centroid: *geojson.NewPointGeometry([]float64{54.60205078125, -25.52509317964987}),
			Value:    2000,
		}
		f3 := models.Block{
			ID:       "F3:0",
			Name:     "FAZENDA 3",
			ParentID: "0",
			Centroid: *geojson.NewPointGeometry([]float64{-355.1165771484375, 52.3755991766591}),
			Value:    3000,
		}
		b0 := models.Block{
			ID:       "B0:F1",
			Name:     "Bloco 0",
			ParentID: "F1",
			Centroid: *geojson.NewPointGeometry([]float64{-354.66064453125, 43.30919109985686}),
			Value:    100,
		}
		b1 := models.Block{
			ID:       "B1:F1",
			Name:     "BLOCO 1",
			ParentID: "F1",
			Centroid: *geojson.NewPointGeometry([]float64{-431.27929687499994, 46.830133640447386}),
			Value:    200,
		}
		b2 := models.Block{
			ID:       "B2:F2",
			Name:     "BLOCO 2",
			ParentID: "F2",
			Centroid: *geojson.NewPointGeometry([]float64{-439.32128906249994, 43.70759350405294}),
			Value:    300,
		}
		b3 := models.Block{
			ID:       "B3:F3",
			Name:     "BLOCO 3",
			ParentID: "F3",
			Centroid: *geojson.NewPointGeometry([]float64{-483.123779296875, 49.25346477497736}),
			Value:    400,
		}
		b4 := models.Block{
			ID:       "B4:F3",
			Name:     "BLOCO 4",
			ParentID: "F3",
			Centroid: *geojson.NewPointGeometry([]float64{-475.16967773437494, 36.16448788632064}),
			Value:    500,
		}
		b5 := models.Block{
			ID:       "B5:F3",
			Name:     "BLOCO 5",
			ParentID: "F3",
			Centroid: *geojson.NewPointGeometry([]float64{-431.971435546875, -13.549881446917126}),
			Value:    600,
		}

		blocks = []models.Block{c0, f1, f2, f3, b0, b1, b2, b3, b4, b5}
	}

	for _, block := range blocks {
		err := db.Set(database.CTX, block.ID, block, 0).Err()
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
