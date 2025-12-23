package repository

import (
	"context"
	"errors"

	"bookhub/internal/domain/entity"
	"bookhub/internal/domain/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const usersCollection = "users"

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) repository.UserRepository {
	return &mongoUserRepository{
		collection: db.Collection(usersCollection),
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *entity.User) error {
	doc := toUserDocument(user)
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *mongoUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return doc.toEntity(), nil
}

func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return doc.toEntity(), nil
}

func (r *mongoUserRepository) List(ctx context.Context, page, limit int) ([]*entity.User, int, error) {
	skip := int64((page - 1) * limit)
	limitInt64 := int64(limit)

	filter := bson.M{}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limitInt64).
		SetSort(bson.D{{Key: "createdat", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []userDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(docs))
	for i, doc := range docs {
		users[i] = doc.toEntity()
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, int(count), nil
}

func (r *mongoUserRepository) Update(ctx context.Context, user *entity.User) error {
	filter := bson.M{"id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"email":     user.Email,
			"active":    user.Active,
			"updatedat": user.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
