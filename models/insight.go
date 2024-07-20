package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Insight struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Text       string             `json:"text"` // TODO: check possible length for string. Might need to change.
	CreatedOn  time.Time          `json:"created_on"`
	ModifiedOn time.Time          `json:"modified_on"`
}

func (i Insight) Description() string {
	return i.Text
}
func (i Insight) GetId() primitive.ObjectID {
	return i.ID
}
func (i Insight) GetAssetType() AssetInterface {
	return i
}
