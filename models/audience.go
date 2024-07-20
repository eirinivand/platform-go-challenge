package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options" // TODO
)

type Audience struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" validate:"required"`
	Gender       Gender             `json:"gender,omitempty"`
	BirthCountry string             `json:"birth_country,omitempty" validate:"country_code"`
	AgeGroups    Range              `json:"age_groups,omitempty"`
	Attributes   []Attribute        `json:"social_commonalities,omitempty"`
}

func (a Audience) Description() string {
	return a.Name
}
func (a Audience) GetId() primitive.ObjectID {
	return a.ID
}
func (a Audience) GetAssetType() AssetInterface {
	return a
}
