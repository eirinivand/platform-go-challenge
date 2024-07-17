package database

import (
	"context"
	"errors"
	"favourites/models"
	"favourites/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options" // TODO
)

type AudienceService interface {
	GetAll(ctx context.Context) ([]models.Audience, error)
	GetByID(ctx context.Context, id string) (models.Audience, error)
	Create(ctx context.Context, m *models.Audience) error
	Update(ctx context.Context, id string, m models.Audience) error
	Delete(ctx context.Context, id string) error
}

type audienceService struct {
	C *mongo.Collection
}

var _ AudienceService = (*audienceService)(nil)

func NewAudienceService(collection *mongo.Collection) AudienceService {

	//indexOpts := new(options.IndexOptions)
	//indexOpts.SetName("audienceIndex").
	//	SetUnique(true)
	//
	//collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	//	Keys:    []string{"_id", "name"},
	//	Options: indexOpts,
	//})

	return &audienceService{C: collection}
}

func (s *audienceService) GetAll(ctx context.Context) ([]models.Audience, error) {
	// Note:
	// The mongodb's go-driver's docs says that you can pass `nil` to "find all" but this gives NilDocument error,
	// probably it's a bug or a documentation's mistake, you have to pass `bson.D{}` instead.
	cur, err := s.C.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Audience

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			return nil, err
		}

		var elem models.Audience
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	return results, nil
}

func (s *audienceService) GetByID(ctx context.Context, id string) (models.Audience, error) {
	var audience models.Audience
	filter, err := utils.MatchID(id)
	if err != nil {
		return audience, err
	}

	err = s.C.FindOne(ctx, filter).Decode(&audience)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return audience, errors.New(utils.ErrorNotFound)
	}
	return audience, err
}

func (s *audienceService) Create(ctx context.Context, m *models.Audience) error {
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	_, err := s.C.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	// The following doesn't work if you have the `bson:"_id` on models.AudienceID field,
	// therefore we manually generate a new ID (look above).
	// res, err := ...InsertOne
	// objectID := res.InsertedID.(primitive.ObjectID)
	// m.ID = objectID
	return nil
}

func (s *audienceService) Update(ctx context.Context, id string, m models.Audience) error {
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

func (s *audienceService) Delete(ctx context.Context, id string) error {
	filter, err := utils.MatchID(id)
	if err != nil {
		return err
	}
	_, err = s.C.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(utils.ErrorNotFound)
		}
		return err
	}

	return nil
}
