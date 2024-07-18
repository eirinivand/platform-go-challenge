package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chart struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string             `json:"title" validate:"required"`
	XAxis Axis               `json:"x_axis" validate:"required"`
	YAxis Axis               `json:"y_axis" validate:"required"`
	ZAxis Axis               `json:"z_axis,omitempty"`
	Data  [][][]float64      `json:"data" validate:"required"`
}

func (ctx *Chart) Description() string {
	return ctx.Title
}
func (ctx *Chart) GetId() primitive.ObjectID {
	return ctx.ID
}
func (ctx *Chart) GetAssetType() AssetInterface {
	return &Chart{}
}
