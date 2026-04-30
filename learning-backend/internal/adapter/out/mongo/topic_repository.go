package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

type topicRepository struct {
	col *mongo.Collection
}

func NewTopicRepository(db *mongo.Database) *topicRepository {
	return &topicRepository{col: db.Collection("topics")}
}

func (r *topicRepository) Create(ctx context.Context, topic *domain.Topic) error {
	result, err := r.col.InsertOne(ctx, topic)
	if err != nil {
		return err
	}
	topic.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *topicRepository) FindAll(ctx context.Context) ([]*domain.Topic, error) {
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topics []*domain.Topic
	if err := cursor.All(ctx, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

func (r *topicRepository) FindByID(ctx context.Context, id string) (*domain.Topic, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var topic domain.Topic
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&topic)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &topic, err
}

func (r *topicRepository) Update(ctx context.Context, topic *domain.Topic) error {
	oid, err := primitive.ObjectIDFromHex(topic.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, topic)
	return err
}

func (r *topicRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
