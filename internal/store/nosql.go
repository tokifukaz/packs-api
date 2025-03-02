package store

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"packs-api/internal/resources"
	"path/filepath"
)

const ordersCollection = "orders"

type NoSQLStore interface {
	// CheckHealth returns the status of the store.
	CheckHealth(ctx context.Context) bool

	// Close terminates any MongoDB connections gracefully.
	Close() error

	CreateOrder(ctx context.Context, order *resources.Order) error
	GetAllOrders(ctx context.Context) ([]*resources.Order, error)
}

// MongoDB represents a MongoDB client.
type MongoDB struct {
	Client *mongo.Client

	DB *mongo.Database
}

// NewMongoDB returns new a MongoDB client.
func NewMongoDB(uri, dbName, certPath string) (*MongoDB, error) {

	clientOptions := options.Client().ApplyURI(uri)

	if certPath != "" {
		c, err := getCustomTLSConfig(certPath)
		if err != nil {
			return nil, err
		}

		clientOptions.SetTLSConfig(c)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	mongoDB := &MongoDB{
		Client: client,
		DB:     client.Database(dbName),
	}

	return mongoDB, nil
}

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := os.ReadFile(filepath.Clean(caFile))

	if err != nil {
		return nil, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return nil, errors.New("failed parsing pem file")
	}

	return tlsConfig, nil
}

// Close terminates any MongoDB connections gracefully.
func (mongoDB *MongoDB) Close() error {
	return mongoDB.Client.Disconnect(context.TODO())
}

// CheckHealth returns the status of the store.
func (mongoDB *MongoDB) CheckHealth(ctx context.Context) bool {
	err := mongoDB.Client.Ping(ctx, readpref.Primary())

	return err == nil
}

func (mongoDB *MongoDB) CreateOrder(ctx context.Context, order *resources.Order) error {
	collection := mongoDB.DB.Collection(ordersCollection)
	_, err := collection.InsertOne(ctx, order)

	return err
}

func (mongoDB *MongoDB) GetAllOrders(ctx context.Context) ([]*resources.Order, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	orders := make([]*resources.Order, 0)

	coll := mongoDB.DB.Collection(ordersCollection)
	cur, err := coll.Aggregate(ctx, mongo.Pipeline{matchStage})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
