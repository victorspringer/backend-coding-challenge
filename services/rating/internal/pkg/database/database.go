package database

import (
	"context"
	"errors"
	"time"

	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/domain"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/log"
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

	// create unique compound index on the "userId" and "movieId" fields
	// this ensures that a user will have only one rating for a movie
	compoundIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "userid", Value: 1},
			{Key: "movieid", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = coll.Indexes().CreateOne(ctx, compoundIndex)
	if err != nil {
		return nil, err
	}

	// create indexes on the "userId" and "movieId" fields
	userIdIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}},
		Options: options.Index(),
	}
	movieIdIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "movieId", Value: 1}},
		Options: options.Index(),
	}

	_, err = coll.Indexes().CreateOne(ctx, userIdIndex)
	if err != nil {
		return nil, err
	}
	_, err = coll.Indexes().CreateOne(ctx, movieIdIndex)
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
func (db *database) Create(ctx context.Context, rating *domain.ValidatedRating) (*domain.Rating, error) {
	if rating.IsValid() {
		filter := bson.M{
			"userId":  rating.Rating.UserID,
			"movieId": rating.Rating.MovieID,
		}

		update := bson.M{
			"$set": rating.Rating,
		}

		updateOptions := options.Update().SetUpsert(true)

		ctx, cancel := context.WithTimeout(ctx, db.timeout)
		defer cancel()

		_, err := db.collection.UpdateOne(ctx, filter, update, updateOptions)
		if err != nil {
			return nil, err
		}

		return &rating.Rating, nil
	}

	return nil, errors.New("invalid rating data")
}

// FindByUserID implements domain.Repository interface's FindByUserID method.
func (db *database) FindByUserID(ctx context.Context, userID string) ([]*domain.Rating, error) {
	filter := bson.D{{Key: "userId", Value: userID}}

	var list []*domain.Rating

	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	cursor, err := db.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var r domain.Rating
		if err = cursor.Decode(&r); err != nil {
			return nil, err
		}
		list = append(list, &r)
	}

	return list, nil
}

// FindByMovieID implements domain.Repository interface's FindByMovieID method.
func (db *database) FindByMovieID(ctx context.Context, movieID string) ([]*domain.Rating, error) {
	filter := bson.D{{Key: "movieId", Value: movieID}}

	var list []*domain.Rating

	ctx, cancel := context.WithTimeout(ctx, db.timeout)
	defer cancel()

	cursor, err := db.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var r domain.Rating
		if err = cursor.Decode(&r); err != nil {
			return nil, err
		}
		list = append(list, &r)
	}

	return list, nil
}
