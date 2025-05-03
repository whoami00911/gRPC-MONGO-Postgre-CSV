package repository

import (
	"context"
	"gRPC-server/internal/domain"
	"gRPC-server/pkg/logger"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBackend struct {
	db     *mongo.Database
	logger *logger.Logger
}

func MongoInit(db *mongo.Database, logger *logger.Logger) *MongoBackend {
	return &MongoBackend{
		db:     db,
		logger: logger,
	}
}

func (m *MongoBackend) Insert(ctx context.Context, product []domain.Product) error {
	client := m.db.Client()
	session, err := client.StartSession()
	if err != nil {
		log.Fatalf("Ошибка создания сессии: %v", err)
	}
	defer session.EndSession(ctx)

	transactionFunc := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		productsInterface := make([]interface{}, len(product))
		for i, v := range product {
			productsInterface[i] = v
		}
		if len(productsInterface) == 0 {
			//m.logger.Error("No products to insert")
			return nil, domain.ErrNoProducts
		}
		_, err := m.db.Collection(viper.GetString("mongo.collection")).InsertMany(ctx, productsInterface)
		if err != nil {
			m.logger.Errorf("Can't Insert in collection: %s", err)
			return nil, err
		}
		return nil, nil
	}
	_, err = session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		// Если транзакция откатилась, здесь будет ошибка
		m.logger.Errorf("Транзакция не удалась: %s", err)
		return err
	}

	return nil
}

func (m *MongoBackend) List(ctx context.Context, sort domain.SortParams) ([]domain.Product, error) {
	var products []domain.Product
	opts := options.Find()
	sortOpts := bson.D{{Key: sort.SortField, Value: sort.SortAsc}}

	opts.SetSort(sortOpts)
	opts.SetSkip(int64(sort.PagingOffset))
	opts.SetLimit(int64(sort.PagingLimit))

	cursor, err := m.db.Collection(viper.GetString("mongo.collection")).Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			m.logger.Errorf("Error decoding product: %s", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (m *MongoBackend) GetByName(ctx context.Context, product domain.Product) (domain.Product, error) {
	var prod domain.Product
	filter := bson.D{{Key: "name", Value: product.Name}}
	result := m.db.Collection(viper.GetString("mongo.collection")).FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		m.logger.Errorf("GetByName err no documents: %s", result.Err())
		return domain.Product{}, result.Err()
	}
	err := result.Decode(&prod)
	if err != nil {
		m.logger.Errorf("GetByName decode error: %s", err)
		return prod, err
	}
	return prod, nil
}

func (m *MongoBackend) UpdateProduct(ctx context.Context, product domain.Product) error {
	//var prod domain.Product
	filter := bson.D{{Key: "name", Value: product.Name}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "price", Value: product.Price}}},
		{Key: "$inc", Value: bson.D{{Key: "changes_count", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "date_of_change", Value: time.Now()}}}}
	_, err := m.db.Collection(viper.GetString("mongo.collection")).UpdateOne(ctx, filter, update)
	if err != nil {
		m.logger.Errorf("Can't update product: %s", err)
		return err
	}
	return nil
}

// доделать удаление
func (m *MongoBackend) DeleteProduct(ctx context.Context, product domain.Product) error {
	filter := bson.D{{Key: "name", Value: product.Name}}
	_, err := m.db.Collection(viper.GetString("mongo.collection")).DeleteOne(ctx, filter)
	if err != nil {
		m.logger.Errorf("Can't delete product: %s", err)
		return err
	}
	return nil
}
