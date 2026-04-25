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

type sleepRepository struct {
	col *mongo.Collection
}

func NewSleepRepository(db *mongo.Database) *sleepRepository {
	return &sleepRepository{col: db.Collection("sleep")}
}

func (r *sleepRepository) Create(ctx context.Context, record *domain.SleepRecord) error {
	result, err := r.col.InsertOne(ctx, record)
	if err != nil {
		return err
	}
	record.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *sleepRepository) FindByID(ctx context.Context, id string) (*domain.SleepRecord, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var record domain.SleepRecord
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &record, err
}

func (r *sleepRepository) FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.SleepRecord, int64, error) {
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

	var records []*domain.SleepRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (r *sleepRepository) Update(ctx context.Context, record *domain.SleepRecord) error {
	oid, err := primitive.ObjectIDFromHex(record.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, record)
	return err
}

func (r *sleepRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
