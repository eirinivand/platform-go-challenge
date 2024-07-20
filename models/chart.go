package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Chart struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string             `json:"title" validate:"required"`
	XAxis      Axis               `json:"x_axis" validate:"required"`
	YAxis      Axis               `json:"y_axis" validate:"required"`
	Points     []Point            `json:"points" validate:"required"`
	CreatedOn  time.Time          `json:"created_on"  bson:"created_on"`
	ModifiedOn time.Time          `json:"modified_on" bson:"modified_on"`
}

func (c Chart) Description() string {
	return c.Title
}
func (c Chart) GetId() primitive.ObjectID {
	return c.ID
}
func (c Chart) GetAssetType() AssetInterface {
	return c
}
