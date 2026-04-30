package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kittitadkp/jeeb-learning/internal/domain"
	"github.com/kittitadkp/jeeb-learning/pkg/apperror"
)

type userRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{col: db.Collection("users")}
}

func (r *userRepository) FindByKeycloakID(ctx context.Context, keycloakID string) (*domain.User, error) {
	var user domain.User
	err := r.col.FindOne(ctx, bson.M{"keycloak_id": keycloakID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &user, err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrNotFound
	}
	var user domain.User
	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, apperror.ErrNotFound
	}
	return &user, err
}

func (r *userRepository) Upsert(ctx context.Context, user *domain.User) error {
	filter := bson.M{"keycloak_id": user.KeycloakID}
	update := bson.M{
		"$set": bson.M{
			"email":        user.Email,
			"display_name": user.DisplayName,
			"updated_at":   time.Now(),
		},
		"$setOnInsert": bson.M{
			"keycloak_id": user.KeycloakID,
			"created_at":  user.CreatedAt,
		},
	}
	opts := options.Update().SetUpsert(true)
	res, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if res.UpsertedID != nil {
		user.ID = res.UpsertedID.(primitive.ObjectID).Hex()
	}
	return nil
}
