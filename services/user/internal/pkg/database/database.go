package database

import (
	"context"
	"fmt"
	"time"

	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/log"
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

	return &database{
		logger:     logger,
		client:     client,
		name:       name,
		collection: client.Database(name).Collection(collection),
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
		ub, err := bson.Marshal(user)
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(ctx, db.timeout)
		defer cancel()

		_, err = db.collection.InsertOne(ctx, ub)
		if err != nil {
			return nil, err
		}

		return &user.User, nil
	}

	return nil, fmt.Errorf(
		"invalid user (id: %s, username: %s, md5 password: %s, name: %s, picture: %s)",
		user.ID, user.Username, user.Password, user.Name, user.Picture,
	)
}

// FindByID implements domain.Repository interface's FindByID method.
func (db *database) FindByID(ctx context.Context, id string) (*domain.User, error) {
	filter := bson.D{{Key: "id", Value: id}}

	var u domain.User

	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	if err := db.collection.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

// FindByUsername implements domain.Repository interface's FindByUsername method.
func (db *database) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.D{{Key: "username", Value: username}}

	var u domain.User

	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	if err := db.collection.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}
