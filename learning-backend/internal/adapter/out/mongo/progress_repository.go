package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

type progressRepository struct {
	col *mongo.Collection
}

func NewProgressRepository(db *mongo.Database) *progressRepository {
	return &progressRepository{col: db.Collection("progress")}
}

func (r *progressRepository) Upsert(ctx context.Context, p *domain.UserProgress) error {
	if p.ID != "" {
		oid, err := primitive.ObjectIDFromHex(p.ID)
		if err != nil {
			return apperror.ErrInternal
		}
		_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, p)
		return err
	}

	filter := bson.M{"user_id": p.UserID, "item_id": p.ItemID}
	update := bson.M{"$set": p}
	opts := options.Update().SetUpsert(true)
	res, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.UpsertedID != nil {
		p.ID = res.UpsertedID.(primitive.ObjectID).Hex()
	}
	return nil
}

func (r *progressRepository) FindByUserAndTopic(ctx context.Context, userID, topicID string) ([]*domain.UserProgress, error) {
	cursor, err := r.col.Find(ctx, bson.M{"user_id": userID, "topic_id": topicID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*domain.UserProgress
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *progressRepository) FindByUserAndItem(ctx context.Context, userID, itemID string) (*domain.UserProgress, error) {
	var p domain.UserProgress
	err := r.col.FindOne(ctx, bson.M{"user_id": userID, "item_id": itemID}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &p, err
}

func (r *progressRepository) DeleteByUserAndTopic(ctx context.Context, userID, topicID string) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"user_id": userID, "topic_id": topicID})
	return err
}
