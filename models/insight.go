package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Insight struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text string             `json:"text"` // TODO: check possible length for string. Might need to change.
}

func (ctx *Insight) Description() string {
	return ctx.Text
}
func (ctx *Insight) GetId() primitive.ObjectID {
	return ctx.ID
}
func (ctx *Insight) GetAssetType() AssetInterface {
	return &Audience{}
}
