package models

import (
	"encoding/json"
	geojson "github.com/paulmach/go.geojson"
)

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

func (b Block) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Block) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, b)
}
