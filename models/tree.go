package models

type Tree struct {
	Block    `json:"block,omitempty"`
	Children []Block `json:"children,omitempty"`
}
