package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	logger     *log.Logger
	client     *mongo.Client
	name       string
	collection *mongo.Collection
	timeout    time.Duration
}

// New returns a new instance of database.
func New(ctx context.Context, logger *log.Logger, uri, name, collection string, timeout time.Duration) (domain.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	coll := client.Database(name).Collection(collection)

	// create unique index on the "id" field
	idIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = coll.Indexes().CreateOne(ctx, idIndex)
	if err != nil {
		return nil, err
	}

	return &database{
		logger:     logger,
		client:     client,
		name:       name,
		collection: coll,
		timeout:    timeout,
	}, nil
}

// Close implements domain.Repository interface's Close method.
func (db *database) Close(ctx context.Context) error {
	if err := db.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

// Create implements domain.Repository interface's Create method.
func (db *database) Create(ctx context.Context, user *domain.ValidatedUser) (*domain.User, error) {
	if user.IsValid() {
		ctx, cancel := context.WithTimeout(ctx, db.timeout)
		defer cancel()

		_, err := db.collection.InsertOne(ctx, user.User)
		if err != nil {
			return nil, err
		}

		return &user.User, nil
	}

	return nil, errors.New("invalid user data")
}

// FindByID implements domain.Repository interface's FindByID method.
func (db *database) FindByID(ctx context.Context, id string) (*domain.User, error) {
	filter := bson.D{{Key: "id", Value: id}}

	var u domain.User

	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	if err := db.collection.FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user %s doesn't exist", id)
		}
		return nil, err
	}

	return &u, nil
}
