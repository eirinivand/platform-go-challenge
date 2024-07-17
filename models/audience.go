package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options" // TODO
)

type Audience struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" validate:"required"`
	Gender       Gender             `json:"gender,omitempty"`
	BirthCountry string             `json:"country,omitempty" validate:"country_code"`
	AgeGroups    Range              `json:"age_groups,omitempty"`
	Attributes   []Attribute        `json:"social_commonalities,omitempty"`
}

func (c *Audience) Description() string {
	return c.Name
}
func (c *Audience) GetId() primitive.ObjectID {
	return c.ID
}
func (c *Audience) GetAssetType() AssetInterface {
	return &Audience{}
}
