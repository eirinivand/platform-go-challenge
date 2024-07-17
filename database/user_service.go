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
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id string) (models.User, error)
	Create(ctx context.Context, m *models.User) error
	CreateAll(ctx context.Context, m []*models.User) error
	Update(ctx context.Context, id string, m models.User) (int64, error)
	Delete(ctx context.Context, id string) error
	GetByUsername(ctx context.Context, username string) (models.User, error)
}

type userService struct {
	C *mongo.Collection
}

var _ UserService = (*userService)(nil)

func NewUserService(collection *mongo.Collection) UserService {
	// indexOpts := new(options.IndexOptions)
	// indexOpts.SetName("userIndex").
	// 	SetUnique(true).
	// 	SetBackground(true).
	// 	SetSparse(true)

	// collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	// 	Keys:    []string{"_id", "name"},
	// 	Options: indexOpts,
	// })

	return &userService{C: collection}
}

func (s *userService) GetAll(ctx context.Context) ([]models.User, error) {
	cur, err := s.C.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.User

	for cur.Next(ctx) {
		if err = cur.Err(); err != nil {
			return nil, err
		}

		//	elem := bson.D{}
		var elem models.User
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		// results = append(results, models.User{ID: elem[0].Value.(primitive.ObjectID)})

		results = append(results, elem)
	}

	return results, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	filter, err := utils.MatchID(id)
	if err != nil {
		return user, err
	}

	err = s.C.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, errors.New(utils.ErrorNotFound)
	}
	return user, err
}

func (s *userService) GetByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	filter := bson.D{{Key: "username", Value: username}}
	err := s.C.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, errors.New(utils.ErrorNotFound)
	}
	return user, err

}

func (s *userService) Create(ctx context.Context, m *models.User) error {
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := errors.New("password encryption  failed")

		return err
	}

	m.Password = string(pass)

	_, err = s.C.InsertOne(context.TODO(), m)

	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = s.C.InsertOne(ctx, m)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// The following doesn't work if you have the `bson:"_id` on models.User.ID field,
	// therefore we manually generate a new ID (look above).
	// res, err := ...InsertOne
	// objectID := res.InsertedID.(primitive.ObjectID)
	// m.ID = objectID
	return nil
}

func (s *userService) CreateAll(ctx context.Context, uAll []*models.User) error {

	var allUsers []interface{}
	for _, i := range uAll {
		allUsers = append(allUsers, i)
	}
	_, err := s.C.InsertMany(context.TODO(), allUsers)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (s *userService) Update(ctx context.Context, id string, m models.User) (int64, error) {
	filter, err := utils.MatchID(id)
	if err != nil {
		return 0, err
	}

	update := bson.D{
		{Key: "$set", Value: m},
	}
	//elem := bson.D{}
	//
	//if m.Title != "" {
	//	elem = append(elem, bson.E{Key: "name", Value: m.Title})
	//}
	//
	//if m.Description != "" {
	//	elem = append(elem, bson.E{Key: "description", Value: m.Description})
	//}
	//
	//if m.Asset != nil {
	//	elem = append(elem, bson.E{Key: "asset", Value: m.Asset})
	//}
	//
	//update := bson.D{
	//	{Key: "$set", Value: elem},
	//}

	cur, err := s.C.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errors.New(utils.ErrorNotFound)
		}
		return 0, err
	}

	return cur.ModifiedCount, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
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
