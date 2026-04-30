package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
	"github.com/kittitadkp/jeeb-learning/pkg/pagination"
)

type itemRepository struct {
	col *mongo.Collection
}

func NewItemRepository(db *mongo.Database) *itemRepository {
	return &itemRepository{col: db.Collection("items")}
}

func (r *itemRepository) Create(ctx context.Context, item *domain.Item) error {
	result, err := r.col.InsertOne(ctx, item)
	if err != nil {
		return err
	}
	item.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *itemRepository) InsertMany(ctx context.Context, items []*domain.Item) error {
	docs := make([]interface{}, len(items))
	for i, item := range items {
		docs[i] = item
	}
	_, err := r.col.InsertMany(ctx, docs)
	return err
}

func (r *itemRepository) FindByTopicID(ctx context.Context, topicID string, opts pagination.Params, category string) ([]*domain.Item, int64, error) {
	filter := bson.M{"topic_id": topicID}
	if category != "" {
		filter["category"] = category
	}

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

	var items []*domain.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *itemRepository) FindAllByTopicID(ctx context.Context, topicID string) ([]*domain.Item, error) {
	cursor, err := r.col.Find(ctx, bson.M{"topic_id": topicID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*domain.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *itemRepository) FindByID(ctx context.Context, id string) (*domain.Item, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var item domain.Item
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &item, err
}

func (r *itemRepository) Update(ctx context.Context, item *domain.Item) error {
	oid, err := primitive.ObjectIDFromHex(item.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, item)
	return err
}

func (r *itemRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
