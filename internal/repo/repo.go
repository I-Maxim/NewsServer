package repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"newsServer/internal/domain"
	"time"
)

type Repository interface {
	Save(context.Context, ...*domain.Article) error
	List(context.Context) ([]*domain.Article, error)
	Load(context.Context, string) (*domain.Article, error)
}

var _ Repository = &MongoRepo{}

type MongoRepo struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

func NewMongoRepo(uri, db, collection string) (Repository, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = client.Ping(ctx, readpref.Primary())
	cancel()
	if err != nil {
		return nil, err
	}

	return &MongoRepo{
		client:         client,
		dbName:         db,
		collectionName: collection,
	}, nil
}

func (mr *MongoRepo) Close(ctx context.Context) error {
	return mr.client.Disconnect(ctx)
}

func (mr *MongoRepo) collection() *mongo.Collection {
	return mr.client.Database(mr.dbName).Collection(mr.collectionName)
}

func (mr *MongoRepo) Save(ctx context.Context, articles ...*domain.Article) error {
	values := make([]interface{}, len(articles))
	for i, v := range articles {
		values[i] = v
	}
	_, err := mr.collection().InsertMany(ctx, values)
	return err
}

func (mr *MongoRepo) List(ctx context.Context) ([]*domain.Article, error) {
	cur, err := mr.collection().Find(ctx, bson.D{}, options.Find().SetProjection(bson.D{{"content", 0}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	result := make([]*domain.Article, 0)
	for cur.Next(ctx) {
		doc := domain.Article{}
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}
		result = append(result, &doc)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (mr *MongoRepo) Load(ctx context.Context, id string) (*domain.Article, error) {
	result := domain.Article{}
	err := mr.collection().FindOne(ctx, bson.D{{"id", id}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}
