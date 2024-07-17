package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Favourite struct {
	ID          primitive.ObjectID `json:"id"         bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	FavouredOn  time.Time          `json:"favoured_on"`
	AssetType   string             `json:"asset_type"  validate:"required,oneof=Chart Insight Audience"`

	//
	AssetId primitive.ObjectID `json:"-"       validate:"required"`

	// This is omitted for bson since no such element exists in DB
	Asset AssetInterface `json:"asset,omitempty" bson:"-"`
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
