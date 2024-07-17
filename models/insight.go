package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Insight struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text string             `json:"text"` // TODO: check possible length for string. Might need to change.
}

func (c *Insight) Description() string {
	return c.Text
}
func (c *Insight) GetId() primitive.ObjectID {
	return c.ID
}
func (c *Insight) GetAssetType() AssetInterface {
	return &Audience{}
}
