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

const booksCollection = "books"

type mongoBookRepository struct {
	collection *mongo.Collection
}

func NewMongoBookRepository(db *mongo.Database) repository.BookRepository {
	return &mongoBookRepository{
		collection: db.Collection(booksCollection),
	}
}

func (r *mongoBookRepository) Create(ctx context.Context, book *entity.Book) error {
	doc := toBookDocument(book)
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *mongoBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	var doc bookDocument
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return doc.toEntity(), nil
}

func (r *mongoBookRepository) GetByISBN(ctx context.Context, isbn string) (*entity.Book, error) {
	var doc bookDocument
	err := r.collection.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return doc.toEntity(), nil
}

func (r *mongoBookRepository) List(ctx context.Context, page, limit int, availableOnly *bool) ([]*entity.Book, int, error) {
	skip := int64((page - 1) * limit)
	limitInt64 := int64(limit)

	filter := bson.M{}
	if availableOnly != nil && *availableOnly {
		filter["availablecopies"] = bson.M{"$gt": 0}
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limitInt64).
		SetSort(bson.D{{Key: "createdat", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []bookDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	books := make([]*entity.Book, len(docs))
	for i, doc := range docs {
		books[i] = doc.toEntity()
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return books, int(count), nil
}

func (r *mongoBookRepository) Update(ctx context.Context, book *entity.Book) error {
	filter := bson.M{"id": book.ID}
	update := bson.M{
		"$set": bson.M{
			"title":           book.Title,
			"author":          book.Author,
			"isbn":            book.ISBN,
			"publishedyear":   book.PublishedYear,
			"totalcopies":     book.TotalCopies,
			"availablecopies": book.AvailableCopies,
			"updatedat":       book.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
