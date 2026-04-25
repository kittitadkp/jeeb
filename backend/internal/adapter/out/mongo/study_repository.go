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

type studyRepository struct {
	col *mongo.Collection
}

func NewStudyRepository(db *mongo.Database) *studyRepository {
	return &studyRepository{col: db.Collection("studies")}
}

func (r *studyRepository) Create(ctx context.Context, session *domain.StudySession) error {
	result, err := r.col.InsertOne(ctx, session)
	if err != nil {
		return err
	}
	session.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *studyRepository) FindByID(ctx context.Context, id string) (*domain.StudySession, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var session domain.StudySession
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &session, err
}

func (r *studyRepository) FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.StudySession, int64, error) {
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

	var sessions []*domain.StudySession
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, 0, err
	}
	return sessions, total, nil
}

func (r *studyRepository) Update(ctx context.Context, session *domain.StudySession) error {
	oid, err := primitive.ObjectIDFromHex(session.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, session)
	return err
}

func (r *studyRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
