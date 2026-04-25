package mongo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kittitad/jeeb/internal/domain"
	"github.com/kittitad/jeeb/pkg/apperror"
	"github.com/kittitad/jeeb/pkg/pagination"
)

type workoutRepository struct {
	col *mongo.Collection
}

func NewWorkoutRepository(db *mongo.Database) *workoutRepository {
	return &workoutRepository{col: db.Collection("workouts")}
}

func (r *workoutRepository) Create(ctx context.Context, workout *domain.Workout) error {
	result, err := r.col.InsertOne(ctx, workout)
	if err != nil {
		return err
	}
	workout.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *workoutRepository) FindByID(ctx context.Context, id string) (*domain.Workout, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var workout domain.Workout
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&workout)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &workout, err
}

func (r *workoutRepository) FindByUserID(ctx context.Context, userID string, opts pagination.Params) ([]*domain.Workout, int64, error) {
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

	var workouts []*domain.Workout
	if err := cursor.All(ctx, &workouts); err != nil {
		return nil, 0, err
	}
	return workouts, total, nil
}

func (r *workoutRepository) Update(ctx context.Context, workout *domain.Workout) error {
	oid, err := primitive.ObjectIDFromHex(workout.ID)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.ReplaceOne(ctx, bson.M{"_id": oid}, workout)
	return err
}

func (r *workoutRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrNotFound
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func parseSortOpt(sort string) (string, int) {
	if strings.HasPrefix(sort, "-") {
		return sort[1:], -1
	}
	return sort, 1
}
