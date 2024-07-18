package database

import (
	"context"
	"errors"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type FavouriteService interface {
	GetAll(ctx context.Context, role string) ([]models.Favourite, error)
	GetByID(ctx context.Context, id string, role string) (models.Favourite, error)
	Create(ctx context.Context, m *models.Favourite) error
	Update(ctx context.Context, id string, m models.Favourite) error
	Delete(ctx context.Context, id string, role string) error
}

type favouriteService struct {
	C *mongo.Collection
}

var _ FavouriteService = (*favouriteService)(nil)

func NewFavouriteService(collection *mongo.Collection) FavouriteService {
	return &favouriteService{C: collection}
}

func (s *favouriteService) GetAll(ctx context.Context, role string) ([]models.Favourite, error) {
	cur, err := s.C.Find(ctx, bson.D{{"role", role}})
	if err != nil {
		panic(err)
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Favourite

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			panic(err)
			return nil, err
		}

		//	elem := bson.D{}
		var elem models.Favourite
		err = cur.Decode(&elem)
		if err != nil {
			panic(err)
			return nil, err
		}

		results = append(results, elem)
	}

	return results, nil
}

func (s *favouriteService) GetByID(ctx context.Context, id string, role string) (models.Favourite, error) {
	var f models.Favourite
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return f, err
	}

	filter := bson.D{{"_id", objectID}, {"role", role}}

	err = s.C.FindOne(context.TODO(), filter).Decode(&f)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return f, errors.New(utils.ErrorNotFound)
	}

	f.EvaluateAssetType()
	fmt.Println(f.GetAssetCollectionByType())

	err = utils.GetDB().Collection(f.GetAssetCollectionByType()).
		FindOne(nil, bson.D{{"_id", f.AssetId}}).Decode(f.Asset)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return f, errors.New(utils.ErrorNotFound)
	}
	return f, nil
}

func (s *favouriteService) Create(ctx context.Context, m *models.Favourite) error {
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}
	m.FavouredOn = time.Now()

	_, err := s.C.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *favouriteService) Update(ctx context.Context, id string, m models.Favourite) error {
	filter, err := utils.MatchID(id)
	if err != nil {
		return err
	}

	update := bson.D{
		{Key: "$set", Value: m},
	}

	_, err = s.C.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(utils.ErrorNotFound)
		}
		return err
	}

	return nil
}

func (s *favouriteService) Delete(ctx context.Context, id string, role string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objectID}, {"role", role}}
	result, err := s.C.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New(utils.ErrorNotFound)
	}
	return nil
}
