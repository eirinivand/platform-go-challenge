package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AssetInterface interface {
	GetId() primitive.ObjectID
	Description() string
}

type AssetCollection struct {
	Charts    []Chart    `json:"charts,omitempty"`
	Insights  []Insight  `json:"insights,omitempty"`
	Audiences []Audience `json:"audiences,omitempty"`
}
