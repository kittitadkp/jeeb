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

type eventRepository struct {
	col *mongo.Collection
}

func NewEventRepository(db *mongo.Database) *eventRepository {
	return &eventRepository{col: db.Collection("events")}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	result, err := r.col.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	event.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *eventRepository) FindByID(ctx context.Context, id string) (*domain.Event, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var event domain.Event
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&event)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &event, err
}

func (r *eventRepository) FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Event, int64, error) {
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

	var events []*domain.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	oid, err := primitive.ObjectIDFromHex(event.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, event)
	return err
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
