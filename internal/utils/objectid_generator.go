package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectIDGenerator interface {
	GenerateRandomObjectID() primitive.ObjectID
	ParseObjectID(id string) (primitive.ObjectID, error)
}

type RandomObjectIDGenerator struct{}

func NewRandomObjectIDGenerator() *RandomObjectIDGenerator {
	return &RandomObjectIDGenerator{}
}

func (og *RandomObjectIDGenerator) GenerateRandomObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (og *RandomObjectIDGenerator) ParseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
