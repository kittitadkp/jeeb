package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/apperror"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type financeRepository struct {
	col *mongo.Collection
}

func NewFinanceRepository(db *mongo.Database) *financeRepository {
	return &financeRepository{col: db.Collection("finance")}
}

func (r *financeRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	result, err := r.col.InsertOne(ctx, tx)
	if err != nil {
		return err
	}
	tx.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *financeRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var tx domain.Transaction
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&tx)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &tx, err
}

func (r *financeRepository) FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Transaction, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	sortField, sortDir := parseSortOpt(opts.Sort)
	findOpts := options.Find().
		SetSkip(opts.Skip()).
		SetLimit(int64(opts.Limit)).
		SetSort(bson.D{{Key: sortField, Value: sortDir}})

	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var txs []*domain.Transaction
	if err := cursor.All(ctx, &txs); err != nil {
		return nil, 0, err
	}
	return txs, total, nil
}

func (r *financeRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	oid, err := primitive.ObjectIDFromHex(tx.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, tx)
	return err
}

func (r *financeRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
