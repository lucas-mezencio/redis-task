package models

import geojson "github.com/paulmach/go.geojson"

//type Tree struct {
//	Block    Block
//	Children []Block
//}

type Block struct {
	ID       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	ParentID string           `json:"parentID,omitempty"`
	Centroid geojson.Geometry `json:"centroid,omitempty"`
	Value    float64          `json:"value,omitempty"`
}
