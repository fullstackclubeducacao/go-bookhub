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

const loansCollection = "loans"

type mongoLoanRepository struct {
	loansCollection *mongo.Collection
	usersCollection *mongo.Collection
	booksCollection *mongo.Collection
}

func NewMongoLoanRepository(db *mongo.Database) repository.LoanRepositoryWithDetails {
	return &mongoLoanRepository{
		loansCollection: db.Collection(loansCollection),
		usersCollection: db.Collection(usersCollection),
		booksCollection: db.Collection(booksCollection),
	}
}

func (r *mongoLoanRepository) Create(ctx context.Context, loan *entity.Loan) error {
	doc := toLoanDocument(loan)
	_, err := r.loansCollection.InsertOne(ctx, doc)
	return err
}

func (r *mongoLoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Loan, error) {
	var doc loanDocument
	err := r.loansCollection.FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return doc.toEntity(), nil
}

func (r *mongoLoanRepository) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*repository.LoanWithDetails, error) {
	loan, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if loan == nil {
		return nil, nil
	}

	userName, err := r.getUserName(ctx, loan.UserID)
	if err != nil {
		return nil, err
	}

	bookTitle, err := r.getBookTitle(ctx, loan.BookID)
	if err != nil {
		return nil, err
	}

	return &repository.LoanWithDetails{
		Loan:      loan,
		UserName:  userName,
		BookTitle: bookTitle,
	}, nil
}

func (r *mongoLoanRepository) GetActiveByUserAndBook(ctx context.Context, userID, bookID uuid.UUID) (*entity.Loan, error) {
	filter := bson.M{
		"userid": userID,
		"bookid": bookID,
		"status": entity.LoanStatusActive,
	}

	var doc loanDocument
	err := r.loansCollection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return doc.toEntity(), nil
}

func (r *mongoLoanRepository) List(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*entity.Loan, int, error) {
	skip := int64((page - 1) * limit)
	limitInt64 := int64(limit)

	filter := r.buildFilter(userID, status)

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limitInt64).
		SetSort(bson.D{{Key: "borrowedat", Value: -1}})

	cursor, err := r.loansCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []loanDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	loans := make([]*entity.Loan, len(docs))
	for i, doc := range docs {
		loans[i] = doc.toEntity()
	}

	count, err := r.loansCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return loans, int(count), nil
}

func (r *mongoLoanRepository) ListWithDetails(ctx context.Context, page, limit int, userID *uuid.UUID, status *string) ([]*repository.LoanWithDetails, int, error) {
	loans, count, err := r.List(ctx, page, limit, userID, status)
	if err != nil {
		return nil, 0, err
	}

	loansWithDetails := make([]*repository.LoanWithDetails, len(loans))
	for i, loan := range loans {
		userName, err := r.getUserName(ctx, loan.UserID)
		if err != nil {
			return nil, 0, err
		}

		bookTitle, err := r.getBookTitle(ctx, loan.BookID)
		if err != nil {
			return nil, 0, err
		}

		loansWithDetails[i] = &repository.LoanWithDetails{
			Loan:      loan,
			UserName:  userName,
			BookTitle: bookTitle,
		}
	}

	return loansWithDetails, count, nil
}

func (r *mongoLoanRepository) Update(ctx context.Context, loan *entity.Loan) error {
	filter := bson.M{"id": loan.ID}
	update := bson.M{
		"$set": bson.M{
			"returnedat": loan.ReturnedAt,
			"status":     loan.Status,
		},
	}

	_, err := r.loansCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoLoanRepository) buildFilter(userID *uuid.UUID, status *string) bson.M {
	filter := bson.M{}

	if userID != nil {
		filter["userid"] = *userID
	}
	if status != nil {
		filter["status"] = *status
	}

	return filter
}

func (r *mongoLoanRepository) getUserName(ctx context.Context, userID uuid.UUID) (string, error) {
	var user struct {
		Name string `bson:"name"`
	}
	err := r.usersCollection.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil
		}
		return "", err
	}
	return user.Name, nil
}

func (r *mongoLoanRepository) getBookTitle(ctx context.Context, bookID uuid.UUID) (string, error) {
	var book struct {
		Title string `bson:"title"`
	}
	err := r.booksCollection.FindOne(ctx, bson.M{"id": bookID}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil
		}
		return "", err
	}
	return book.Title, nil
}
