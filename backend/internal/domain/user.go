package domain

import "time"

type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	KeycloakID  string    `bson:"keycloak_id" json:"keycloak_id"`
	Email       string    `bson:"email" json:"email"`
	DisplayName string    `bson:"display_name" json:"display_name"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}
