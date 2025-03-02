package resources

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRequest struct {
	Items     int   `json:"items"`
	PackSizes []int `json:"packSizes"`
}

type Order struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Items        int                `json:"items" bson:"items"`
	PackSizes    []int              `json:"packSizes" bson:"pack_sizes"`
	PackQuantity map[int]int        `json:"packQuantity" bson:"pack_quantity"`
	CreatedAt    time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updated_at"`
}
