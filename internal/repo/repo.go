package repo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"newsServer/internal/domain"
	"time"
)

type Repository interface {
	Save(context.Context, ...*domain.ArticleDB) error
	List(context.Context) ([]*domain.ArticleDB, error)
	Load(context.Context, string) (*domain.ArticleDB, error)
}

var _ Repository = &MongoRepo{}

type MongoRepo struct {
	client         *mongo.Client
	cancel         context.CancelFunc
	dbName         string
	collectionName string
}

func NewMongoRepo(uri, db, collection string) (Repository, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancelCtx()
		return nil, fmt.Errorf("connect to mongodb failed: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err = client.Ping(ctx, readpref.Primary())
	cancel()
	if err != nil {
		cancelCtx()
		return nil, fmt.Errorf("ping mongodb failed: %w", err)
	}

	return &MongoRepo{
		client:         client,
		cancel:         cancelCtx,
		dbName:         db,
		collectionName: collection,
	}, nil
}

func (mr *MongoRepo) Close(ctx context.Context) error {
	mr.cancel()
	return mr.client.Disconnect(ctx)
}

func (mr *MongoRepo) collection() *mongo.Collection {
	return mr.client.Database(mr.dbName).Collection(mr.collectionName)
}

func (mr *MongoRepo) Save(ctx context.Context, articles ...*domain.ArticleDB) error {
	values := make([]interface{}, len(articles))
	for i, v := range articles {
		values[i] = v
	}
	_, err := mr.collection().InsertMany(ctx, values)
	return err
}

func (mr *MongoRepo) List(ctx context.Context) ([]*domain.ArticleDB, error) {
	cur, err := mr.collection().Find(ctx, bson.D{}, options.Find().SetProjection(bson.D{{Key: "content", Value: 0}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	result := make([]*domain.ArticleDB, 0)
	for cur.Next(ctx) {
		doc := domain.ArticleDB{}
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

func (mr *MongoRepo) Load(ctx context.Context, id string) (*domain.ArticleDB, error) {
	result := domain.ArticleDB{}
	err := mr.collection().FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
