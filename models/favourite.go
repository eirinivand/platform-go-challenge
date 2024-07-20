package models

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Favourite struct {
	ID          primitive.ObjectID `json:"id"             bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	FavouredOn  time.Time          `json:"favoured_on"`
	AssetType   string             `json:"asset_type"  validate:"required,oneof=Chart Insight Audience"`
	AssetId     primitive.ObjectID `json:"asset_id"    validate:"required"`
	Asset       AssetInterface     `json:"asset"`
	Role        string             `json:"-"           validate:"required"`
}

// Make sure these match the types of assets that exist
// And more importantly the Favourite.AssetType validation.oneof list
const (
	CHART_ASSET         = "Chart"
	INSIGHT_ASSET       = "Insight"
	AUDIENCE_ASSET      = "Audience"
	CHART_COLLECTION    = "charts"
	INSIGHT_COLLECTION  = "insights"
	AUDIENCE_COLLECTION = "audiences"
)

func (f *Favourite) GetAssetCollectionByType() string {

	switch f.AssetType {
	case CHART_ASSET:
		return CHART_COLLECTION
	case INSIGHT_ASSET:
		return INSIGHT_COLLECTION
	case AUDIENCE_ASSET:
		return AUDIENCE_COLLECTION
	default:
		panic("invalid asset type" + f.AssetType)
	}
	return ""
}

func (f *Favourite) EvaluateAssetType() {

	switch f.AssetType {
	case CHART_ASSET:
		f.Asset = new(Chart)
	case INSIGHT_ASSET:
		f.Asset = new(Insight)
	case AUDIENCE_ASSET:
		f.Asset = new(Audience)
	default:
		panic("invalid asset type" + f.AssetType)
	}
}

func (f *Favourite) UnmarshalBSON(data []byte) error {

	var raw map[string]interface{}
	err := bson.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	fmt.Println(raw)
	var ok bool
	f.AssetType, ok = raw["assettype"].(string)
	fmt.Println(f.AssetType, ok)
	f.ID, ok = raw["_id"].(primitive.ObjectID)
	fmt.Println(f.ID, ok)
	f.Title, ok = raw["title"].(string)
	fmt.Println(f.Title, ok)
	f.AssetId, ok = raw["assetid"].(primitive.ObjectID)
	fmt.Println(f.AssetId, ok)
	f.Role, ok = raw["role"].(string)
	fmt.Println(f.Role, ok)
	fOn, ok := raw["favouredon"].(primitive.DateTime)
	f.FavouredOn = fOn.Time()
	fmt.Println(f.FavouredOn, ok)
	f.Description, _ = raw["description"].(string)
	fmt.Println(f.Description, ok)
	assetBytes, err := bson.Marshal(raw["asset"])
	if err != nil {
		return err
	}
	switch f.AssetType {
	case CHART_ASSET:
		var a Chart
		err = bson.Unmarshal(assetBytes, &a)
		if err != nil {
			return err
		}
		f.Asset = a
	case INSIGHT_ASSET:
		var a Insight
		err = bson.Unmarshal(assetBytes, &a)
		if err != nil {
			return err
		}
		f.Asset = a
	case AUDIENCE_ASSET:
		var a Audience
		err = bson.Unmarshal(assetBytes, &a)
		if err != nil {
			return err
		}
		f.Asset = a
	default:
		return errors.New("invalid asset type" + f.AssetType)
	}

	return nil

}
